Natural Language Processing Workflow
====

## Types

### Word

- raw value
- cleaned value

### Token
- name
- raw value
- confidence

## Sentence Processing Workflow

- detect language
- start tokenization by generating words with the splitter
  - splitter dictionaries on raw data first
  - sentence splitter afterwards
    - remove punctuation
    - split whitespaces
    - words
- raw words
  - create raw tokens
  - save raw format
  - to lower case
  - remove diacritics
- filter stop words
- lemmatize
- tokenize using dictionaries comparaison with levenstein
