resource "github_repository" "repo" {
  name                   = var.name
  description            = var.description
  homepage_url           = var.homepage_url
  private                = var.private
  has_issues             = var.hasIssues
  has_projects           = var.HasProjects
  has_wiki               = var.hasWiki
  allow_merge_commit     = var.allow_merge_commit
  allow_squash_merge     = var.allow_squash_merge
  allow_rebase_merge     = var.allow_rebase_merge
  has_downloads          = var.hasDownloads
  delete_branch_on_merge = var.delete_branch_on_merge
  default_branch         = var.default_branch
  archived               = var.archived
  topics                 = var.topics
}
