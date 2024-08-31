resource "github_repository" "repo2" {
  name      = "repo2"
  auto_init = true

  visibility  = "visibility2"
  description = "a description2"

  template {
    owner      = "owner2"
    repository = "repository2"
  }

  topics       = ["topic4", "topic5"]
  homepage_url = "http://localhost/2"

  pages {
    source {
      branch = "branch2"
      path   = "path2"
    }
  }

  has_issues    = true
  has_projects  = true
  has_wiki      = true
  has_downloads = true

  allow_merge_commit     = true
  allow_rebase_merge     = true
  allow_squash_merge     = true
  allow_auto_merge       = true
  delete_branch_on_merge = true

  merge_commit_title          = "aMergeCommitTitle2"
  merge_commit_message        = "aMergeCommitMessage2"
  squash_merge_commit_title   = "aSquashMergeCommitTitle2"
  squash_merge_commit_message = "aSquashMergeCommitMessage2"

  vulnerability_alerts = false

  archived           = true
  archive_on_destroy = false
}

resource "github_branch_default" "repo2" {
  repository = github_repository.repo2.name
  branch     = "master2"
}

resource "github_branch" "repo2-feature-branch2" {
  repository = github_repository.repo2.name
  branch     = "feature/branch2"

  lifecycle {
    ignore_changes = [source_branch]
  }
}

resource "github_branch" "repo2-feature-branch3" {
  repository    = github_repository.repo2.name
  branch        = "feature/branch3"
  source_branch = "branch3-source-branch2"
  source_sha    = "branch3-source-sha2"
}

resource "github_branch_protection" "repo2-default" {
  repository_id           = github_repository.repo2.node_id
  pattern                 = github_branch_default.repo2.branch
  enforce_admins          = true
  allows_deletions        = true
  allows_force_pushes     = true
  push_restrictions       = ["default-branch-pushRestriction2"]
  required_linear_history = true
  require_signed_commits  = true

  required_status_checks {
    strict   = true
    contexts = ["default-branch-context2"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = true
    restrict_dismissals             = true
    dismissal_restrictions          = ["default-branch-dismissalRestriction2"]
    require_code_owner_reviews      = true
    required_approving_review_count = 1
  }
}

resource "github_branch_protection" "repo2-feature_SLASH_branch2" {
  repository_id           = github_repository.repo2.node_id
  pattern                 = "feature/branch2"
  enforce_admins          = false
  allows_deletions        = false
  allows_force_pushes     = false
  push_restrictions       = ["branch2-pushRestriction2"]
  required_linear_history = false
  require_signed_commits  = false

  required_status_checks {
    strict   = false
    contexts = ["branch2-context2"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = false
    restrict_dismissals             = false
    dismissal_restrictions          = ["branch2-dismissalRestriction2"]
    require_code_owner_reviews      = false
    required_approving_review_count = 2
  }
}

resource "github_branch_protection" "repo2-feature_SLASH_branch3" {
  repository_id           = github_repository.repo2.node_id
  pattern                 = "feature/branch3"
  enforce_admins          = true
  allows_deletions        = true
  allows_force_pushes     = true
  push_restrictions       = ["branch3-pushRestriction2"]
  required_linear_history = true
  require_signed_commits  = true

  required_status_checks {
    strict   = true
    contexts = ["branch3-context2"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = true
    restrict_dismissals             = true
    dismissal_restrictions          = ["branch3-dismissalRestriction2"]
    require_code_owner_reviews      = true
    required_approving_review_count = 3
  }
}

resource "github_branch_protection" "repo2-a-pattern2" {
  repository_id           = github_repository.repo2.node_id
  pattern                 = "a-pattern2"
  enforce_admins          = false
  allows_deletions        = false
  allows_force_pushes     = false
  push_restrictions       = ["branch-protection-pushRestriction2"]
  required_linear_history = false
  require_signed_commits  = false

  required_status_checks {
    strict   = false
    contexts = ["branch-protection-context2"]
  }

  required_pull_request_reviews {
    dismiss_stale_reviews           = false
    restrict_dismissals             = false
    dismissal_restrictions          = ["branch-protection-dismissalRestriction2"]
    require_code_owner_reviews      = false
    required_approving_review_count = 4
  }
}
