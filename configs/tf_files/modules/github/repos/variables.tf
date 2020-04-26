variable "name" {
  type = string
}
variable "description" {
  type    = string
  default = ""
}
variable "homepage_url" {
  type    = string
  default = ""
}
variable "private" {
  type    = bool
  default = true
}
variable "hasIssues" {
  type    = bool
  default = true
}
variable "HasProjects" {
  type    = bool
  default = false
}
variable "hasWiki" {
  type    = bool
  default = true
}
variable "allow_merge_commit" {
  type    = bool
  default = true
}
variable "allow_squash_merge" {
  type    = bool
  default = true
}
variable "allow_rebase_merge" {
  type    = bool
  default = true
}
variable "hasDownloads" {
  type    = bool
  default = true
}
variable "delete_branch_on_merge" {
  type    = bool
  default = true
}
variable "default_branch" {
  type    = string
  default = "master"
}
variable "archived" {
  type    = bool
  default = false
}
variable "topics" {
  type    = list(string)
  default = []
}

