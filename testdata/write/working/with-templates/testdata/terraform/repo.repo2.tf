resource "github_repository" "repo2" {
  name = "repo2"

  has_issues    = true
  has_projects  = true
  has_wiki      = true
  has_downloads = true

  allow_merge_commit     = false
  allow_rebase_merge     = false
  allow_squash_merge     = false
  delete_branch_on_merge = false

  vulnerability_alerts = true
}
