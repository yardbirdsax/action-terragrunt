terraform {
  source = "${get_repo_root()}/integration/module"
}

inputs = {
  input = "hello"
}