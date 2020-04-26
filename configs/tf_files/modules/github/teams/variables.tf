variable "team_name" {
  type = string
}
variable "team_description" {
  type    = string
  default = ""
}
variable "team_privacy" {
  type    = string
  default = "secret"
}
variable "team_members" {
  type = map(string)
}

