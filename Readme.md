# Compare messages files

## Check missing

- print missing messages which are in `PATH_TO_SOURCE_CORE_MSG_FILE` but missing in `PATH_TO_DST_TRANSLATED_MSG_FILE`

```bash
./compare-msg.exe -f1 PATH_TO_SOURCE_CORE_MSG_FILE  -f2 PATH_TO_DST_TRANSLATED_MSG_FILE
```

## Translation file

CSV file with contains lines:

`NAMESPACE.MESSAGE_ID,ORIGINAL_MESSAGE,TRANSLATES_MESSAGE`

## Check missing with translation file

- print missing messages which are in `SOURCE_CORE_MSG_FILE` but missing in `DST_TRANSLATED_MSG_FILE` and also missing in `TRANSLATION_FILE`

```bash
./compare-msg.exe -f1 PATH_TO_SOURCE_CORE_MSG_FILE  -f2 PATH_TO_DST_TRANSLATED_MSG_FILE -translation-file PATH_TO_TRANSLATION_FILE
```

## Add missing translations

- add missing translation from TRANSLATION_FILE and save it to PATH_TO_OUT_DST_TRANSLATED_MSG_FILE 

```bash
./compare-msg.exe -f1 PATH_TO_SOURCE_CORE_MSG_FILE  -f2 PATH_TO_DST_TRANSLATED_MSG_FILE -translation-file PATH_TO_TRANSLATION_FILE -out PATH_TO_OUT_DST_TRANSLATED_MSG_FILE -overwrite
```
