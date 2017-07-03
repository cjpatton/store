/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 *
 * dict.h -- A probabilistic dictionary structure designed by Denis Charles
 * and Kumar Chellapilla ("Bloomier Filters: A second look").
 *
 * Represents a key/value store. The keys are of arbitrary length, and the
 * values are of at most some fixed length. It is assumed that each of the keys
 * are unique.
 *
 * The construction uses a simple, acyclic graph, and so this program includes
 * some basic structures for representing and maninpulating undirected graphs.
 */

#ifndef DICT_H
#define DICT_H

#include "tiny.h"

#define MAX_OUT_DEGREE 32
#define NODE_CT_FACTOR 2.09

// Stores a node, including some state associated to it used for various
// algorithms.
typedef struct {

  // Adjacent nodes in the graph are indicated by their indices in the adjacency
  // list representing the graph.
  int adj_node [MAX_OUT_DEGREE];

  // Label of each edge connecting this node to its neighbor, e.g., an integer
  // weight. In our case, this is an index in a list of strings.
  int adj_edge [MAX_OUT_DEGREE];

  int adj_ct,  // Total number of outgoing edges
      rec,     // State for cycle checking
      visited; // State for depth-first search
} node_t;

// Stores a graph as an adjacency list.
typedef struct {
  node_t *node;
  int node_ct;
} graph_t;

// Returns a pointer to a new graph with node_ct nodes.
//
// Must be freed with graph_free().
graph_t *graph_new(int node_ct);

// Frees memory allocated to graph.
void graph_free(graph_t *graph);

// Checks if the input graph is simple (meaning no multi-edges or self-loops)
// and acyclic.
//
// Assumes the graph is undirected, meaning for each edge (x,y), there is an
// edge (y,x). This is not counted as a cycle.
//
// Returns:
//  - OK if the graph is simple and acyclic,
//  - ERR_GRAPH_CYCLE if graph has a cycle or self-loop, and
//  - ERR_GRAPH_MULTIEDGE if graph has a multi-edge.
int graph_simple_and_acyclic(graph_t *graph);

// Stores parameters for dictionaries, including the salt.
typedef struct {

  int table_length,    // Number of table rows
      max_value_bytes, // Maximum length of any allowed value
      tag_bytes,       // Number of bytes used for tag
      row_bytes,       // max_value_bytes+tag_bytes+1
      salt_bytes;      // Length of the salt

  char *salt; // Of length salt_bytes + 1.

} dict_params_t;

// Stores the dictionary structure.
typedef struct {
  dict_params_t params;
  char *table; // Of length row_bytes * table_length
} dict_t;

// Stores a compressed dictionary, which reduces space at the cost of O(log
// item_ct) query time instead of O(1).
typedef struct {
  dict_params_t params;
  int compressed_table_length;
  char *table;
  int *idx; // Index of each row in uncompressed table.
} cdict_t;

// Computes optimal table length.
//
// Optimal table length is ceil(item_ct * 2.09.
int dict_compute_table_length(int item_ct);

// Returns a pointer to a new dict_t
//
// Lets row_bytes = max_value_bytes + tag_bytes + 1. The extra byte is for
// padding the value. Must be freed with dict_free(). Returns NULL if row_bytes
// > HASH_BYTES (defined in const.h).
dict_t *dict_new(int table_length, int max_value_bytes, int tag_bytes,
    int salt_bytes);

// Frees memory allocated to dict.
void dict_free(dict_t *dict);

// Generates a fresh seed and constructs the "hash graph" from sequence 'key' of
// length 'item_ct'.
//
// Returns:
//  - OK if success,
//  - ERR_GRAPH_EXCEEDED_MAX_OUT_DEGREE if some node has more than
//    MAX_OUT_DEGREE edges,
//  - ERR_HMAC or ERR_SHA512 if tiny->h erred erred, or
//  - BAD if tiny->h failed to find a uniform output.
int dict_generate_graph(dict_t *ctx, tiny_ctx *tiny, char **key,
    int *key_bytes, int item_ct, graph_t *graph);

// Constructs a dictionary from the sequences 'key' and 'value', each of length
// 'item_ct'.
//
// Returns:
//  - OK if successful,
//  - ERR_DICT_TOO_MANY_ITEMS,
//  - ERR_DICT_LONG_VALUE, or
//  - ERR_HMAC.
int dict_create(dict_t *dict, tiny_ctx *tiny, char **key, int *key_bytes,
    char **value, int *value_bytes, int item_ct);

// Computes row indices in dict->table corresponding to 'key' and sets *x and *y
// to these values.
//
// Returns:
//  - OK if successful
//  - ERR_HMAC or ERR_SHA512 if tiny->h erred, or
//  - BAD if tinyprf or tinyhash failed to find a uniform output.
int dict_compute_rows(dict_params_t params, tiny_ctx *tiny, const char *key,
    int key_bytes, int *x, int *y);

// Computes value from from 'key' and rows 'xrow' and 'yrow' in the table.
// Writes output to 'value' and the number of bytes to '*value_bytes'.
//
// Returns:
//  - OK if successful,
//  - ERR_DICT_BAD_KEY if key was not found,
//  - ERR_DICT_BAD_PADDING if the padding was incorrect, or
//  - ERR_HMAC or ERR_SHA512 if tiny->g errs.
int dict_compute_value(dict_params_t params, tiny_ctx *tiny, const char *key,
    int key_bytes, const char *xrow, const char *yrow, char *value,
    int *value_bytes);

// Checks if 'key' is stored by 'dict'. If so, set 'value' to the value
// associated to 'key' and '*value_bytes' to the length of 'value'.
//
// Returns OK if successful; otherwise, this function propagates errors from
// dict_compute_rows() and dict_compute_value().
int dict_get(dict_t *dict, tiny_ctx *tiny, const char *key, int key_bytes,
    char *value, int *value_bytes);

// Computes a new cdict_t* from a dict_t* and returns it.
//
// Destroy with cdict_free().
cdict_t *dict_compress(dict_t *dict);

// Computes index in compressed table associated to x, an output of
// dict_compute_value().
//
// Initial values: l=0 and r=comp->compressed_table_length.
int cdict_binsearch(cdict_t *comp, int x, int l, int r);

// Free memory allocated for 'comp'.
void cdict_free(cdict_t *comp);

// Checks if 'key' is stored by 'dict'. If so, set 'value' to the value
// associated to 'key' and '*value_bytes' to the length of 'value'.
//
// Returns OK if successful; otherwise, this function propagates errors from
// dict_compute_rows() and dict_compute_value().
int cdict_get(cdict_t *comp, tiny_ctx *tiny,
    const char *key, int key_bytes, char *value, int *value_bytes);

#endif //DICT_H
