# -*-ruby-*-
listen = ":8080"
bodylimit = 10000

apicontext = [
  {
    context = "root"
    key = "0fddbff5d8010f81cd28a7d77f3e38981b13d6164c2fd6e1c3f60a4287630c37",
  },
  {
    context = "foo",
    key = "970b391f22f515d96b3e9b86a2c62c627968828e47b356994d2e583188b4190a"
  }
]

#url = "https://sokrates.daemon.de"

# this is the root context with all permissions
super = "root"

mail = {
  server = "localhost"
  port = "25"
  from = "root@localhost"
  password = ""
}
