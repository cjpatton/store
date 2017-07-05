/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 */
#include "const.h"
#include "dict.h"

#include "assert.h"
#include "openssl/rand.h"
#include "stdio.h"
#include "stdlib.h"
#include "string.h"

#define PAD_BYTE 0x70
#define ZERO_BYTE 0x00
#define ROW(row) (row*(dict->params.row_bytes))

graph_t *graph_new(int node_ct) {
  graph_t *graph = malloc(sizeof(graph_t));
  graph->node = malloc(node_ct * sizeof(node_t));
  graph->node_ct = node_ct;
  return graph;
}

void graph_free(graph_t *graph) {
  free(graph->node);
  free(graph);
}

// Recursively checks that component associated with graph->node[i] is simple
// and acyclic.
//
// Called by graph_simple_and_acyclic().
int node_simple_and_acyclic(graph_t *graph, int x, int rec) {
  int err, y;

  if (!graph->node[x].visited) {
    graph->node[x].visited = 1;
    graph->node[x].rec = rec;

    for (int j = 0; j < graph->node[x].adj_ct; j++) {

      y = graph->node[x].adj_node[j];

      // Check for self-loop.
      if (x == y) {
          return ERR_GRAPH_CYCLE;
      }

      // Check for multi-edge.
      for (int k = 0; k < graph->node[x].adj_ct; k++) {
        if (k != j && graph->node[x].adj_node[k] == y) {
          return ERR_GRAPH_MULTI_EDGE;
        }
      }

      // Check for cycle.
      if (graph->node[y].visited && rec - graph->node[y].rec > 1) {
        return ERR_GRAPH_CYCLE;
      }

      // Recursively check child.
      err = node_simple_and_acyclic(graph, y, rec+1);
      if (err != OK) {
        return err;
      }
    }
  }

  return OK;
}

int graph_simple_and_acyclic(graph_t *graph) {
  int err = OK;
  for (int i = 0; i < graph->node_ct; i++) {
    if (!graph->node[i].visited) {
      err = node_simple_and_acyclic(graph, i, 0);
      if (err != OK) {
        break;
      }
    }
  }
  for (int i = 0; i < graph->node_ct; i++) {
    graph->node[i].rec = 0;
    graph->node[i].visited = 0;
  }
  return err;
}

///////////////////////////////////////////////////////////////////////////////

int dict_compute_table_length(int item_ct) {
  int table_length = ((double)item_ct * NODE_CT_FACTOR) * 10;
  if ((table_length % 10) == 0) {
    return table_length / 10;
  } else {
    return (table_length / 10) + 1;
  }
}

dict_t *dict_new(int table_length, int max_value_bytes, int tag_bytes,
    int salt_bytes) {
  if ((max_value_bytes + tag_bytes + 1) > HASH_BYTES) {
    return NULL;
  }
  dict_t *dict = malloc(sizeof(dict_t));
  dict->params.table_length = table_length;
  dict->params.max_value_bytes = max_value_bytes;
  dict->params.tag_bytes = tag_bytes;
  dict->params.row_bytes = tag_bytes + max_value_bytes + 1;
  dict->params.salt_bytes = salt_bytes;
  dict->table = malloc(table_length * (dict->params.row_bytes) * sizeof(char));
  dict->params.salt = malloc((salt_bytes + 1) * sizeof(char));
  return dict;
}

void dict_free(dict_t *dict) {
  free(dict->table);
  free(dict->params.salt);
  free(dict);
}

int dict_compute_rows(dict_params_t params, tiny_ctx *tiny, const char *key, int
    key_bytes, int *x, int *y) {

  if (tiny->_use_prf) {
    params.salt[params.salt_bytes] = 1;
    *x = tinyprf(tiny, key, key_bytes, params.salt, params.salt_bytes+1);
    if (*x < 0) {
      return *x;
    }

    params.salt[params.salt_bytes] = 2;
    *y = tinyprf(tiny, key, key_bytes, params.salt, params.salt_bytes+1);
    if (*y < 0) {
      return *y;
    }
  } else {
    params.salt[params.salt_bytes] = 1;
    *x = tinyhash(tiny, key, key_bytes, params.salt, params.salt_bytes+1);
    if (*x < 0) {
      return *x;
    }

    params.salt[params.salt_bytes] = 2;
    *y = tinyhash(tiny, key, key_bytes, params.salt, params.salt_bytes+1);
    if (*y < 0) {
      return *y;
    }
  }

  return OK;
}


int dict_generate_graph(dict_t *dict, tiny_ctx *tiny, char **key,
    int *key_bytes, int item_ct, graph_t *graph) {

  RAND_bytes((unsigned char *)dict->params.salt, dict->params.salt_bytes);
  for (int i = 0; i < graph->node_ct; i++) {
    graph->node[i].adj_ct = 0;
    graph->node[i].rec = 0;
    graph->node[i].visited = 0;
  }

  for (int j = 0; j < item_ct; j++) {

    // Map key[j] to a pair of node (x,y).
    int x, y;
    int err = dict_compute_rows(dict->params, tiny, key[j], key_bytes[j], &x, &y);
    if (err != OK) {
      return err;
    }

    // Add edge (x,y) to graph.
    int adj_ct = graph->node[x].adj_ct;
    if (adj_ct == MAX_OUT_DEGREE) {
      return ERR_GRAPH_EXCEEDED_MAX_OUT_DEGREE;
    }
    graph->node[x].adj_node[adj_ct] = y;
    graph->node[x].adj_edge[adj_ct] = j;
    graph->node[x].adj_ct ++;

    // Add edge (y,x) to graph.
    adj_ct = graph->node[y].adj_ct;
    if (adj_ct == MAX_OUT_DEGREE) {
      return ERR_GRAPH_EXCEEDED_MAX_OUT_DEGREE;
    }
    graph->node[y].adj_node[adj_ct] = x;
    graph->node[y].adj_edge[adj_ct] = j;
    graph->node[y].adj_ct ++;
  }

  return OK;
}

// First pass of each tree of the graph: For each child y of x, let (in, out) be
// the key/value pair associated to edge (x,y). Compute
//
//    Z = hash(salt || 3 || in) ^ out
//
// and set the row associated with y be Z. Add Z into the row associated to p,
// where p is the root of tree.
int dict_create_traverse1(dict_t *dict, tiny_ctx *tiny, graph_t *graph, int x,
    int p, char **key, int *key_bytes, char **value, int *value_bytes)
{
  int e, y, err, prow = ROW(p);
  graph->node[x].rec = 1;

  // For each child y of x.
  for (int i = 0; i < graph->node[x].adj_ct; i++) {

    y = graph->node[x].adj_node[i];
    e = graph->node[x].adj_edge[i];

    if (graph->node[y].rec == 1) {
      continue;
    }

    // Compute row y.
    dict->params.salt[dict->params.salt_bytes] = 3;
    if (tiny->_use_prf) {
      int err = prf(tiny, key[e], key_bytes[e], dict->params.salt,
                    dict->params.salt_bytes+1, &dict->table[ROW(y)],
                    dict->params.row_bytes);
      if (err != OK) {
        return err;
      }
    } else {
      int err = hash(tiny, key[e], key_bytes[e], dict->params.salt,
                     dict->params.salt_bytes+1, &dict->table[ROW(y)],
                     dict->params.row_bytes);
      if (err != OK) {
        return err;
      }
    }

    // Add value[e] to y.
    int yrow = ROW(y);
    for (int j = 0; j < value_bytes[e]; j++) {
      dict->table[yrow+j] ^= value[e][j];
    }
    dict->table[yrow+value_bytes[e]] ^= PAD_BYTE;

    // Add y to p.
    for (int j = 0; j < dict->params.row_bytes; j++) {
      dict->table[prow+j] ^= dict->table[yrow+j];
    }

    // Recurse on child.
    err = dict_create_traverse1(dict, tiny, graph, y, p, key, key_bytes, value,
        value_bytes);
    if (err != OK) {
      return OK;
    }
  }
  return OK;
}

// Second pass of each tree of the graph: For each child y of x, add row
// associated to the root p into the row associated to y.
void dict_create_traverse2(dict_t *dict, tiny_ctx *tiny, graph_t *graph, int x,
    int p, char **key, int *key_bytes, char **value, int *value_bytes)
{
  int y;
  graph->node[x].rec = 2;

  // For each child y of x.
  for (int i = 0; i < graph->node[x].adj_ct; i++) {
    y = graph->node[x].adj_node[i];

    if (graph->node[y].rec == 2) {
      continue;
    }

    // Add x to y.
    int xrow = ROW(x), yrow = ROW(y);
    for (int j = 0; j < dict->params.row_bytes; j++) {
      dict->table[yrow+j] ^= dict->table[xrow+j];
    }

    // Recurse on child.
    dict_create_traverse2(dict, tiny, graph, y, p, key, key_bytes, value,
        value_bytes);
  }
}

// TODO Exit do-while loop if the number of retries exceeds 1000.
int dict_create(dict_t *dict, tiny_ctx *tiny, char **key, int *key_bytes,
    char **value, int *value_bytes, int item_ct) {

  if (item_ct >= dict->params.table_length) {
    return ERR_DICT_TOO_MANY_ITEMS;
  }
  for (int i = 0; i < item_ct; i++) {
    if (value_bytes[i] > dict->params.max_value_bytes) {
      return ERR_DICT_LONG_VALUE;
    }
  }

  // Clear the table.
  memset(dict->table, 0, dict->params.table_length * (dict->params.row_bytes));

  // Generate a random, simple, and acyclic graph.
  graph_t *graph = graph_new(dict->params.table_length);
  int err, ct = 0;
  do {
    err = dict_generate_graph(dict, tiny, key, key_bytes, item_ct, graph);
    if (err == OK) {
      err = graph_simple_and_acyclic(graph);
    }
    ct ++;
  } while (err != OK);

  // Compute table.
  //
  // Assumes graph->node[x].rec == 0 for all x and dict->table is
  // zeroed-out.
  for (int x = 0; x < graph->node_ct; x++) {
    if (graph->node[x].rec == 0) {
      err = dict_create_traverse1(
          dict, tiny, graph, x, x, key, key_bytes, value, value_bytes);
      if (err != OK) {
        break;
      }
      dict_create_traverse2(
          dict, tiny, graph, x, x, key, key_bytes, value, value_bytes);
    }
  }

  graph_free(graph);
  return err;
}

///////////////////////////////////////////////////////////////////////////////

int dict_compute_value(dict_params_t params, tiny_ctx *tiny, const char *key, int
    key_bytes, const char *xrow, const char *yrow, char *value, int *value_bytes) {

  // Compute pad.
  char buf [HASH_BYTES];
  params.salt[params.salt_bytes] = 3;
  if (tiny->_use_prf) {
    int err = prf(tiny, key, key_bytes, params.salt, params.salt_bytes+1, buf, HASH_BYTES);
    if (err != OK) {
      return err;
    }
  } else {
    int err = hash(tiny, key, key_bytes, params.salt, params.salt_bytes+1, buf, HASH_BYTES);
    if (err != OK) {
      return err;
    }
  }

  for (int j = 0; j < params.row_bytes; j++) {
    buf[j] ^= xrow[j] ^ yrow[j];
  }

  // Check the tag.
  char t = ZERO_BYTE;
  for (int j = params.max_value_bytes+1; j < params.row_bytes; j++) {
    t |= buf[j];
  }

  // If tag bytes are all 0, then copy value to 'value'.
  if (t == ZERO_BYTE) {
    int last_byte = params.max_value_bytes;
    while (last_byte > 0 && buf[last_byte] == ZERO_BYTE) {
      last_byte--;
    }
    if (buf[last_byte] != PAD_BYTE) {
      return ERR_DICT_BAD_PADDING;
    }
    memcpy(value, buf, last_byte);
    *value_bytes = last_byte;
  } else {
    return ERR_DICT_BAD_KEY;
  }

  return OK;
}

int dict_get(dict_t *dict, tiny_ctx *tiny, const char *key, int key_bytes,
    char *value, int *value_bytes) {

  int x, y;
  int err = dict_compute_rows(dict->params, tiny, key, key_bytes, &x, &y);
  if (err != OK) {
    return err;
  }

  char *xrow = &dict->table[x * dict->params.row_bytes],
       *yrow = &dict->table[y * dict->params.row_bytes];

  return dict_compute_value(dict->params, tiny, key, key_bytes, xrow, yrow, value,
      value_bytes);
}

///////////////////////////////////////////////////////////////////////////////

cdict_t *dict_compress(dict_t *dict) {
  int salt_bytes = dict->params.salt_bytes,
      row_bytes = dict->params.row_bytes;
  cdict_t *compressed = malloc(sizeof(cdict_t));
  compressed->params.table_length = dict->params.table_length;
  compressed->params.max_value_bytes = dict->params.max_value_bytes;
  compressed->params.tag_bytes = dict->params.tag_bytes;
  compressed->params.row_bytes = row_bytes;
  compressed->params.salt_bytes = salt_bytes;

  // Copy salt.
  compressed->params.salt = malloc(salt_bytes + 1);
  memcpy(compressed->params.salt, dict->params.salt, salt_bytes);
  compressed->params.salt[salt_bytes] = ZERO_BYTE;

  // Count The non-zero rows.
  int row_ct = 0;
  for (int x = 0; x < dict->params.table_length; x++) {
    int xrow = row_bytes * x;
    char c = ZERO_BYTE;
    for (int j = 0; j < row_bytes; j++) {
      c |= dict->table[j+xrow];
    }
    row_ct += (c != ZERO_BYTE);
  }
  compressed->compressed_table_length = row_ct;

  // Create table and idx.
  compressed->table = malloc((row_ct+1) * row_bytes);
  compressed->idx = malloc(row_ct * sizeof(int));
  int y = 0;
  for (int x = 0; x < dict->params.table_length; x++) {
    int xrow = row_bytes * x;
    int yrow = row_bytes * y;
    char c = ZERO_BYTE;
    for (int j = 0; j < row_bytes; j++) {
      c |= dict->table[j+xrow];
    }
    if (c != ZERO_BYTE) {
      memcpy(&compressed->table[yrow], &dict->table[xrow], row_bytes);
      compressed->idx[y] = x;
      y++;
    }
  }
  memset(&compressed->table[row_ct*row_bytes], 0, row_bytes);

  return compressed;
}

void cdict_free(cdict_t *compressed) {
  free(compressed->idx);
  free(compressed->table);
  free(compressed->params.salt);
  free(compressed);
}

// Get row index by binary search.
int cdict_binsearch(cdict_t *comp, int x, int l, int r) {
  if (r == l) {
    return comp->compressed_table_length;
  }
  int m = (r + l) / 2;
  if (x > comp->idx[m]) {
    return cdict_binsearch(comp, x, m+1, r);
  } else if (x == comp->idx[m]) {
    return m;
  } else {
    return cdict_binsearch(comp, x, l, m);
  }
}

int cdict_get(cdict_t *comp, tiny_ctx *tiny,
    const char *key, int key_bytes, char *value, int *value_bytes) {

  int x, y;
  int err = dict_compute_rows(comp->params, tiny, key, key_bytes, &x, &y);
  if (err != OK) {
    return err;
  }

  // Look up rows in idx.
  int xidx = cdict_binsearch(comp, x, 0, comp->compressed_table_length),
      yidx = cdict_binsearch(comp, y, 0, comp->compressed_table_length);

  char *xrow = &comp->table[xidx * comp->params.row_bytes],
       *yrow = &comp->table[yidx * comp->params.row_bytes];

  return dict_compute_value(comp->params, tiny, key, key_bytes, xrow, yrow,
      value, value_bytes);
}
