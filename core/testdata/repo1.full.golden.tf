resource "github_repository" "repo1" {
  name      = "repo1"
  auto_init = false

  visibility  = "visibility1"
  description = "a description1"

  template {
    owner      = "owner1"
    repository = "repository1"
  }

  topics       = ["topic2", "topic3"]
  homepage_url = "http://localhost/1"

  pages {
    source {
      branch = "branch1"
      path   = "path1"
    }
  }

  has_issues    = false
  has_projects  = false
  has_wiki      = false
  has_downloads = false

  allow_merge_commit     = false
  allow_rebase_merge     = false
  allow_squash_merge     = false
  allow_auto_merge       = false
  delete_branch_on_merge = false

  merge_commit_title          = "aMergeCommitTitle1"
  merge_commit_message        = "aMergeCommitMessage1"
  squash_merge_commit_title   = "aSquashMergeCommitTitle1"
  squash_merge_commit_message = "aSquashMergeCommitMessage1"

  vulnerability_alerts = true

  archived           = false
  archive_on_destroy = true
}

resource "github_branch_default" "repo1" {
  repository = github_repository.repo1.name
  branch     = "master1"
}

resource "github_branch" "repo1-feature-branch1" {
  repository = github_repository.repo1.name
  branch     = "feature/branch1"

  lifecycle {
    ignore_changes = [source_branch]
  }
}

resource "github_branch" "repo1-feature-branch2" {
  repository    = github_repository.repo1.name
  branch        = "feature/branch2"
  source_branch = "branch2-source-branch1"
  source_sha    = "branch2-source-sha1"
}

resource "github_branch_protection" "repo1-default" {
  repository_id           = github_repository.repo1.node_id
  pattern                 = github_branch_default.repo1.branch
  enforce_admins          = false
  allows_deletions        = false
  allows_force_pushes     = false
  push_restrictions       = ["default-branch-pushRestriction1"]
  required_linear_history = false
  require_signed_commits  = false

  required_status_checks {
    strict   = false
    contexts = ["default-branch-context1"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = false
    restrict_dismissals             = false
    dismissal_restrictions          = ["default-branch-dismissalRestriction1"]
    require_code_owner_reviews      = false
    required_approving_review_count = 4
  }
}

resource "github_branch_protection" "repo1-feature_SLASH_branch1" {
  repository_id           = github_repository.repo1.node_id
  pattern                 = "feature/branch1"
  enforce_admins          = true
  allows_deletions        = true
  allows_force_pushes     = true
  push_restrictions       = ["branch1-pushRestriction1"]
  required_linear_history = true
  require_signed_commits  = true

  required_status_checks {
    strict   = true
    contexts = ["branch1-context1"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = true
    restrict_dismissals             = true
    dismissal_restrictions          = ["branch1-dismissalRestriction1"]
    require_code_owner_reviews      = true
    required_approving_review_count = 5
  }
}

resource "github_branch_protection" "repo1-feature_SLASH_branch2" {
  repository_id           = github_repository.repo1.node_id
  pattern                 = "feature/branch2"
  enforce_admins          = false
  allows_deletions        = false
  allows_force_pushes     = false
  push_restrictions       = ["branch2-pushRestriction1"]
  required_linear_history = false
  require_signed_commits  = false

  required_status_checks {
    strict   = false
    contexts = ["branch2-context1"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = false
    restrict_dismissals             = false
    dismissal_restrictions          = ["branch2-dismissalRestriction1"]
    require_code_owner_reviews      = false
    required_approving_review_count = 6
  }
}

resource "github_branch_protection" "repo1-a-pattern1" {
  repository_id           = github_repository.repo1.node_id
  pattern                 = "a-pattern1"
  enforce_admins          = true
  allows_deletions        = true
  allows_force_pushes     = true
  push_restrictions       = ["branch-protection-pushRestriction1"]
  required_linear_history = true
  require_signed_commits  = true

  required_status_checks {
    strict   = true
    contexts = ["branch-protection-context1"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = true
    restrict_dismissals             = true
    dismissal_restrictions          = ["branch-protection-dismissalRestriction1"]
    require_code_owner_reviews      = true
    required_approving_review_count = 0
  }
}
