/* Copyright (c) 2017, Christopher Patton
 * All rights reserved.
 */
#include "stdio.h"
#include "stdlib.h"
#include "string.h"

#include "const.h"
#include "dict.h"
#include "tiny.h"

void print_adj_list(graph_t *graph) {
  for (int i = 0; i < graph->node_ct; i++) {
    printf("%d: ", i);
    for (int j = 0; j < graph->node[i].adj_ct; j++) {
      printf("%d[%d] ", graph->node[i].adj_node[j], graph->node[i].adj_edge[j]);
    }
    printf("\n");
  }
}

void print_table(dict_t *dict) {
  for (int i = 0; i < dict->params.table_length; i++) {
    for (int j = 0; j < dict->params.row_bytes; j++) {
      printf("%02x", (unsigned char)dict->table[j+(i*dict->params.row_bytes)]);
    }
    printf("\n");
  }
}

void print_cdict(cdict_t *comp) {
  printf("salt ");
  for (int j = 0; j < comp->params.salt_bytes; j++) {
    printf("%02x", (unsigned char)comp->params.salt[j]);
  }
  printf("\n");
  for (int x = 0; x < comp->compressed_table_length; x++) {
    int xrow = x * comp->params.row_bytes;
    printf("%-5d", comp->idx[x]);
    for (int j = 0; j < comp->params.row_bytes; j++) {
      printf("%02x", (unsigned char)comp->table[j+xrow]);
    }
    printf("\n");
  }
}

char *key [] = {
  "This",
  "is",
  "cool!"
};

int key_bytes [] = { 4, 2, 5 };

char *value [] = {
  "secret",
  "stuff",
  ""
};

int value_bytes [] = { 6, 5, 0 };

char *value_nopad [] = {
  "01234567",
  "89abcdef",
  "ghijklmn"
};

int value_nopad_bytes [] = { 8, 8, 8 };

int test_dict_new_free() {
  int ret = OK;
  const char test[] = "dict_new_free";

  // Test good dict_new().
  dict_t *dict = dict_new(1000, 16, 3, 8, 1);
  ASSERT_FATAL("dict is NULL", dict != NULL);
  ASSERT_EQ_ERROR("table_length", dict->params.table_length, 1000);
  ASSERT_EQ_ERROR("max_value_bytes", dict->params.max_value_bytes, 16);
  ASSERT_EQ_ERROR("tag_bytes", dict->params.tag_bytes, 3);
  ASSERT_EQ_ERROR("row_bytes", dict->params.row_bytes, 16 + 3 + 1);
  ASSERT_EQ_ERROR("salt_bytes", dict->params.salt_bytes, 8);
  dict_free(dict);

  // Test bad dict_new().
  dict = dict_new(1000, HASH_BYTES, 1, 16, 1);
  ASSERT_FATAL("dict not NULL", dict == NULL);

  return ret;
}

int test_bad_create() {
  int ret = OK;
  const char test[] = "bad_create";

  dict_t *dict = dict_new(2, 2, 3, 8, 1);
  tiny_ctx *tiny = tinyhash_new(2);

  int err = dict_create(dict, tiny, key, key_bytes, value, value_bytes, 3);
  ASSERT_EQ_ERROR("dict_create(3)", err, ERR_DICT_TOO_MANY_ITEMS);

  err = dict_create(dict, tiny, key, key_bytes, value, value_bytes, 2);
  ASSERT_EQ_ERROR("dict_create(2)", err, ERR_DICT_TOO_MANY_ITEMS);

  err = dict_create(dict, tiny, key, key_bytes, value, value_bytes, 1);
  ASSERT_EQ_ERROR("dict_create(1)", err, ERR_DICT_LONG_VALUE);

  tinyhash_free(tiny);
  dict_free(dict);

  return ret;
}

int test_dict_create_get() {
  int ret = OK;
  const char test[] = "dict_create_get";

  int tag_bytes = 2,
      salt_bytes = 8,
      item_ct = 3;

  // Make sure that the value[0] is the longest.
  int max_value_bytes = value_bytes[0];
  int table_length = dict_compute_table_length(item_ct);

  tiny_ctx *tiny = tinyprf_new(table_length);
  tinyprf_init_generate_key(tiny);

  dict_t *dict = dict_new(table_length, max_value_bytes, tag_bytes, salt_bytes, 1);

  char out [max_value_bytes];
  int out_bytes;

  int err, trials = 100;
  for (int trial = 0; trial < trials; trial++) {
    err = dict_create(dict, tiny, key, key_bytes, value, value_bytes, item_ct);
    ASSERT_OK_FATAL("dict_create()", err);
    // FIXME Callback
    //
    // tinyprf_free(tiny);
    // dict_free(dict);

    for (int it = 0; it < item_ct; it++) {
      err = dict_get(dict, tiny, key[it], key_bytes[it], out, &out_bytes);
      ASSERT_OK_ERROR("dict_get()", err);
      ASSERT_EQ_ERROR("dict_get(): out_bytes", out_bytes, value_bytes[it]);
      ASSERT_EQ_ERROR("dict_get(): strncmp(value[it], out, out_bytes)",
        strncmp(value[it], out, out_bytes), 0);
    }
  }

  if (ret != OK) {
    tinyprf_free(tiny);
    dict_free(dict);
    return ret;
  }

  err = dict_get(dict, tiny, "hella", strlen("hella"), out, &out_bytes);
  ASSERT_EQ_ERROR("dict_get()", err, ERR_DICT_BAD_KEY);

  tinyprf_free(tiny);
  dict_free(dict);
  return ret;
}

int test_dict_create_nopad_get() {
  int ret = OK;
  const char test[] = "dict_create_nopad_get";

  int tag_bytes = 2,
      salt_bytes = 8,
      item_ct = 3;

  // Make sure that the value[0] is the longest.
  int max_value_bytes = value_nopad_bytes[0];
  int table_length = dict_compute_table_length(item_ct);

  tiny_ctx *tiny = tinyprf_new(table_length);
  tinyprf_init_generate_key(tiny);

  dict_t *dict = dict_new(table_length, max_value_bytes, tag_bytes, salt_bytes, 0);

  char out [max_value_bytes];
  int out_bytes;

  int err = dict_create(
      dict, tiny, key, key_bytes, value_nopad, value_nopad_bytes, item_ct);
  ASSERT_OK_FATAL("dict_create()", err);
  // FIXME Callback
  //
  // tinyprf_free(tiny);
  // dict_free(dict);

  for (int it = 0; it < item_ct; it++) {
    err = dict_get(dict, tiny, key[it], key_bytes[it], out, &out_bytes);
    ASSERT_OK_ERROR("dict_get()", err);
    ASSERT_EQ_ERROR("dict_get(): out_bytes", out_bytes, value_nopad_bytes[it]);
    ASSERT_EQ_ERROR("dict_get(): strncmp(value[it], out, out_bytes)",
      strncmp(value_nopad[it], out, out_bytes), 0);
  }

  if (ret != OK) {
    tinyprf_free(tiny);
    dict_free(dict);
    return ret;
  }

  err = dict_get(dict, tiny, "hella", strlen("hella"), out, &out_bytes);
  ASSERT_EQ_ERROR("dict_get()", err, ERR_DICT_BAD_KEY);

  tinyprf_free(tiny);
  dict_free(dict);
  return ret;
}

int test_many() {
  int ret = OK;
  const char test[] = "dict_new_free";

  int item_ct = 1000,
      key_ct = 6000,
      max_value_bytes = 8,
      tag_bytes = 2,
      salt_bytes = 8;

  char **key = malloc(key_ct * sizeof(char *));
  int *key_bytes = malloc(key_ct * sizeof(int));
  char **value = malloc(item_ct * sizeof(char *));
  int *value_bytes = malloc(item_ct * sizeof(int));

  for (int i = 0; i < key_ct; i++) {
    key[i] = malloc(max_value_bytes+1);
    sprintf(key[i], "key#%d", i);
    key_bytes[i] = strlen(key[i]);
  }
  for (int i = 0; i < item_ct; i++) {
    value[i] = malloc(max_value_bytes+1);
    sprintf(value[i], "val#%d", i);
    value_bytes[i] = strlen(value[i]);
  }

  int table_length = dict_compute_table_length(item_ct);
  tiny_ctx *tiny = tinyprf_new(table_length);
  tinyprf_init_generate_key(tiny);

  dict_t *dict = dict_new(table_length, max_value_bytes, tag_bytes,
      salt_bytes, 1);

  int err = dict_create(dict, tiny, key, key_bytes, value, value_bytes,
      item_ct);
  ASSERT_OK_FATAL("dict_create()", err);
  // FIXME callback

  char out [max_value_bytes];
  int out_bytes;

  // These should all succeed.
  for (int i = 0; i < item_ct; i++) {
    err = dict_get(dict, tiny, key[i], key_bytes[i], out, &out_bytes);
    if (err != OK) {
      fprintf(stderr, "error: %s: dict_get(key[\"%s\"]) errs\n", test, key[i]);
    }
    ASSERT_EQ_ERROR("dict_get(): strncmp(value[i], out, out_bytes)",
      strncmp(value[i], out, out_bytes), 0);
  }

  // Most of these should fail.
  int bad_ct = 0;
  for (int i = item_ct; i < key_ct; i++) {
    err = dict_get(dict, tiny, key[i], key_bytes[i], out, &out_bytes);
    if (err == OK) {
      bad_ct++;
    }
  }
  if (bad_ct > 0) {
    fprintf(stderr, "warning: %s: bad_ct = %d of %d\n", test, bad_ct,
        key_ct-item_ct);
  }

  for (int i = 0; i < item_ct; i++) {
    free(value[i]);
  }
  for (int i = 0; i < key_ct; i++) {
    free(key[i]);
  }
  free(key);
  free(key_bytes);
  free(value);
  free(value_bytes);

  dict_free(dict);
  tinyprf_free(tiny);

  return ret;
}

int test_cdict() {
  int ret = OK;
  const char test[] = "cdict";

  int item_ct = 20,
      max_value_bytes = 8,
      tag_bytes = 2,
      salt_bytes = 8;

  char **key = malloc(item_ct * sizeof(char *));
  int *key_bytes = malloc(item_ct * sizeof(int));
  char **value = malloc(item_ct * sizeof(char *));
  int *value_bytes = malloc(item_ct * sizeof(int));

  for (int i = 0; i < item_ct; i++) {
    key[i] = malloc(max_value_bytes+1);
    sprintf(key[i], "key#%d", i);
    key_bytes[i] = strlen(key[i]);
    value[i] = malloc(max_value_bytes+1);
    sprintf(value[i], "val#%d", i);
    value_bytes[i] = strlen(value[i]);
  }

  int table_length = dict_compute_table_length(item_ct);
  tiny_ctx *tiny = tinyprf_new(table_length);
  tinyprf_init_generate_key(tiny);

  dict_t *dict = dict_new(table_length, max_value_bytes, tag_bytes,
      salt_bytes, 1);

  int err = dict_create(dict, tiny, key, key_bytes, value, value_bytes,
      item_ct);
  ASSERT_OK_FATAL("dict_create()", err); // FIXME callback

  cdict_t *compressed = dict_compress(dict);
  ASSERT_FATAL("dict_compress(dict) = NULL", compressed != NULL); // FIXME callback
  ASSERT_EQ_ERROR("uncompressed_table_length",
      compressed->params.table_length, dict->params.table_length);
  ASSERT_EQ_ERROR("max_value_bytes", compressed->params.max_value_bytes,
      dict->params.max_value_bytes);
  ASSERT_EQ_ERROR("tag_bytes", compressed->params.tag_bytes, dict->params.tag_bytes);
  ASSERT_EQ_ERROR("row_bytes", compressed->params.row_bytes, dict->params.row_bytes);
  ASSERT_EQ_ERROR("salt_bytes", compressed->params.salt_bytes, dict->params.salt_bytes);

  ASSERT_EQ_ERROR("strncmp(compressed->salt, dict->params.salt, dict->params.salt_bytes)",
      strncmp(compressed->params.salt, dict->params.salt, dict->params.salt_bytes), 0);

  print_cdict(compressed);

  char out [max_value_bytes];
  int out_bytes;

  for (int i = 0; i < item_ct; i++) {
    err = cdict_get(compressed, tiny, key[i], key_bytes[i], out,
        &out_bytes);
    if (err != OK) {
      fprintf(stderr, "error: %s: cdict_get(key[\"%s\"]) errs\n",
          test, key[i]);
      ret = err;
    } else {
      ASSERT_EQ_ERROR("cdict_get(): strncmp(value[i], out, out_bytes)",
        strncmp(value[i], out, out_bytes), 0);
    }
  }

  for (int i = 0; i < item_ct; i++) {
    free(key[i]);
    free(value[i]);
  }
  free(key);
  free(key_bytes);
  free(value);
  free(value_bytes);

  cdict_free(compressed);
  dict_free(dict);
  tinyprf_free(tiny);

  return ret;
}

int main() {
  if (test_dict_new_free()==OK &&
      test_bad_create()==OK &&
      test_dict_create_get()==OK &&
      test_dict_create_nopad_get()==OK &&
      test_many()==OK &&
      test_cdict()==OK) {
    printf("pass\n");
    return 0;
  } else {
    return 1;
  }
}
