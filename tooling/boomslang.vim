
syntax match bsText /text\zs.*$/

syntax match bsKeyword /is/
syntax match bsKeyword /the/
syntax match bsKeyword /of/
syntax match bsKeyword /if/
syntax match bsKeyword /otif/
syntax match bsKeyword /otherwise/
syntax match bsKeyword /ask/
syntax match bsKeyword /while/
syntax match bsKeyword /for/
syntax match bsKeyword /break/
" syntax match bsKeyword /text/

syntax match bsBuiltin /show/
syntax match bsBuiltin /debug/
syntax match bsBuiltin /plus/
syntax match bsBuiltin /minus/
syntax match bsBuiltin /multiply/
syntax match bsBuiltin /divide/
syntax match bsBuiltin /biggerthan/
syntax match bsBuiltin /smallerthan/
syntax match bsBuiltin /equals/
syntax match bsBuiltin /notequals/


hi def link bsText String
hi def link bsKeyword Keyword
hi def link bsBuiltin Identifier
