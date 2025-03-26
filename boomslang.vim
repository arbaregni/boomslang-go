
syntax match bsKeyword /is/
syntax match bsKeyword /the/
syntax match bsKeyword /of/
syntax match bsKeyword /if/
syntax match bsKeyword /otif/
syntax match bsKeyword /otherwise/

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

syntax match bsText /text\zs.*$/

hi def link bsKeyword Keyword
hi def link bsBuiltin Identifier
hi def link bsText String
