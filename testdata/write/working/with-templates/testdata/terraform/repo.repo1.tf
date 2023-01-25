resource "github_repository" "repo1" {
  name = "repo1"

  description = "a description"

  has_issues    = false
  has_projects  = true
  has_wiki      = true
  has_downloads = true

  allow_merge_commit     = false
  allow_rebase_merge     = false
  allow_squash_merge     = false
  delete_branch_on_merge = false

  archived           = true
  archive_on_destroy = true
}

resource "github_branch_default" "repo1" {
  repository = "repo1"
  branch     = "master"
}

resource "github_branch_protection" "repo1-master" {
  repository_id           = "repo1"
  pattern                 = "master"
  enforce_admins          = true
  allows_deletions        = true
  allows_force_pushes     = true
  push_restrictions       = ["pushRestriction"]
  required_linear_history = true
  require_signed_commits  = true

  required_status_checks {
    strict   = true
    contexts = ["context1"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = false
    restrict_dismissals             = true
    dismissal_restrictions          = ["dismissalRestriction"]
    require_code_owner_reviews      = true
    required_approving_review_count = 2
  }
}
