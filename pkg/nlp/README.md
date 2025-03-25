Simple Natural Language Processing
====

## Features

- support for multiple languages
- support for multiple scripts
- all steps of the processing workflow are configurable and based on regular-expressions
- splitting
- cleaning
- normalisation
- lexeme based tokenisation
- dictionary based tokenisation

## Structure

```
Tokenizer
   |_ LanguageCodes
   |_ Splitter
   |_ Dictionaries   # are applied on the normalised words
   |_ Lexemes        # are applied on the raw sentence
```

## Sources

- https://www.geeksforgeeks.org/token-patterns-and-lexems/
- https://www.geeksforgeeks.org/dictionary-based-tokenization-in-nlp/


## Libs

- https://stackoverflow.com/a/3614928 (Lexer/Parser)
- https://github.com/aaaton/golem (lemmatizer)
- https://github.com/go-ego/gse (text segmentation)
- https://github.com/sentencizer/sentencizer (Sentence Splitting)
- https://github.com/golang-nlp/stopwords/tree/main/data (stopwords)
- https://github.com/gorgonia/gorgonia (machine learning)
- https://github.com/sjwhitworth/golearn (machine learning)
- https://github.com/jdkato/prose (archive NLP lib)
- https://github.com/pemistahl/lingua-go (language detection)
