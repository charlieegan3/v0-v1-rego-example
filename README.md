# v0-v1-rego-example

This repo contains a simple program that shows how to use V0 and V1 Rego in the
same program. The program will output:

```
Rego Version: v1
[{"expressions":[{"value":["foo"],"text":"data.example.messages","location":{"row":1,"col":1}}]}]
Rego Version: v0
[{"expressions":[{"value":["foo"],"text":"data.example.messages","location":{"row":1,"col":1}}]}]
-----------------
[{"expressions":[{"value":["bar","foo"],"text":"data.v1.messages | data.v0.messages","location":{"row":1,"col":1}}]}]
```
