resource "github_membership" "user" {
  username = var.login
  role     = var.role
}