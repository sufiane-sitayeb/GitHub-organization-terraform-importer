resource "github_team" "team" {
  name        = var.team_name
  description = var.team_description
  privacy     = var.team_privacy
}
resource "github_team_membership" "team_membership" {
  for_each = var.team_members

  team_id  = github_team.team.id
  username = each.key
  role     = each.value
}
