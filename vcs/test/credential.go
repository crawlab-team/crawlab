package test

type Credential struct {
	Username               string `json:"username"`
	Password               string `json:"password"`
	TestRepoHttpUrl        string `json:"test_repo_http_url"`
	TestRepoMultiBranchUrl string `json:"test_repo_multi_branch_url"`
	SshUsername            string `json:"ssh_username"`
	SshPassword            string `json:"ssh_password"`
	TestRepoSshUrl         string `json:"test_repo_ssh_url"`
	PrivateKey             string `json:"private_key"`
	PrivateKeyPath         string `json:"private_key_path"`
}
