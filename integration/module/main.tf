provider "local" {

}
resource "local_file" "name" {
  content = "hello world"
  filename = "out.txt"
}