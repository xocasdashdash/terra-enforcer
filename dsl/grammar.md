

# Grammar

digit               =   "0".."9"
ascii_letter        =   "A".."Z" | "a" .. "z"
letter              =   asci_letter | "_"
word                =   (letter){letter|digit}

WithStatement       =   "with" ("ALL"   |   "NONE"  |   "SOME"  |   "ONLY")    .
IDStatement         =   word  { "."   word   }
ValueStatement      =   "[" { word "," } word . "]"
AttributeStatement  =   "attribute" IDStatement WithStatement   ValueStatement
BlockStatement      =   "{" (AttributeStatement){"," AttributeStatement } "}"
ResourceStatement   =   "resource" IDStatement  "has"   BlockStatement  .
Program             =   {   ResourceStatement  }   .

