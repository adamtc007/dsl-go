package ebnf

const Text = `
request = "(" "onboarding-request" meta orchestrator [catalog] ")" .
meta = "(" ":meta" "(" "request-id" String ")" "(" "version" Number ")" [ "(" "created-at" String ")" ] [ "(" "updated-at" String ")" ] ")" .
orchestrator = "(" ":orchestrator" lifecycle entities [resources] [flows] [policies] [product-service-mappings] ")" .
lifecycle = "(" ":lifecycle" "(" "states" Ident* ")" "(" "initial" Ident ")" "(" "transitions" transition* ")" ")" .
transition = "(" "->" Ident Ident [guard] [effects] ")" .
guard = "(" "when" expr ")" .
effects = "(" "do" action-call* ")" .
entities = "(" ":entities" entity* ")" .
entity = "(" "entity" ":id" String ":type" Ident "(" "attrs" attr* ")" ")" .
attr = "(" Ident value [ ":" "provenance" String ] [ ":" "needed-by" "(" Ident* ")" ] ")" .
resources = "(" ":resources" resource* ")" .
resource = "(" "resource" ":id" String ":type" Ident [requires] [config] ")" .
requires = "(" "requires" require-item* ")" .
require-item = "(" "entity" String ")" .
config = "(" "config" kv-pair* ")" .
flows = "(" ":flows" flow* ")" .
flow = "(" "flow" ":id" String [String] "(" "steps" step* ")" ")" .
step = task | gate | fork | join .
task = "(" "task" ":id" String ":on" String ":op" Ident "(" "args" kv-pair* ")" [ "(" "needs" String* ")" ] [ "(" "produces" String* ")" ] [ "(" "labels" Ident* ")" ] ")" .
gate = "(" "gate" ":id" String "(" "when" String ")" ")" .
fork = "(" "fork" ":id" String "(" "branches" String* ")" ")" .
join = "(" "join" ":id" String "(" "after" String* ")" ")" .
policies = "(" ":policies" policy* ")" .
policy = "(" "policy" Ident kv-pair* ")" .
catalog = "(" ":catalog" "(" ":attributes" attr-def* ")" "(" ":actions" action-def* ")" ")" .
attr-def = "(" Ident ":" "type" Ident [ ":" "enum" "(" Ident* ")" ] [ ":" "format" Ident ] [ ":" "pii" ("true" | "false") "] ")" .
action-def = "(" Ident "(" "params" param-def* ")" "(" "needs" String* ")" "(" "produces" String* ")" ")" .
param-def = "(" Ident ":" "type" Ident [ ":" "required" ("true" | "false") ] [ ":" "enum" "(" Ident* ")" ] ")" .
expr = Ident [String] .
kv-pair = "(" Ident value ")" .
value = String | Number | "true" | "false" | Ident .
product-service-mappings = "(" ":product-service-mappings" mapping* ")" .
mapping = "(" "mapping" ":product" String ":services" "(" String* ")" ":resources" "(" String* ")" ")" .

String = \"\" ( { all unicode characters | \\ ( \" \" | \\ ) } ) \"\" .
Number = [ "-" ] { "0" ... "9" } [ "." { "0" ... "9" } ] .
Ident = ( "a" ... "z" | "A" ... "Z" | "_" ) { "a" ... "z" | "A" ... "Z" | "0" ... "9" | "_" | "-" } .
`
