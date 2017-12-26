resource "ns1_team" "t" {
  name = "terraform acc test team %s"
}

resource "ns1_user" "u" {
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"
  teams = ["${ns1_team.t.id}"]
  notify {
    billing = true
  }
}