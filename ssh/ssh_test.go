package ssh

import (
	glog "github.com/andriyg76/glog"
	"github.com/andriyg76/scm-backup/lists"
	"github.com/andriyg76/scm-backup/os"
	"github.com/stretchr/testify/assert"
	os2 "os"
	"testing"
)

func init() {
	glog.SetLevel(glog.TRACE)
}
func TestParseSshAgentOutput(t *testing.T) {
	out := lists.String("SSH_AUTH_SOCK=/tmp/ssh-kQN4kvgauzrv/agent.4465; export SSH_AUTH_SOCK;",
		"SSH_AGENT_PID=4466; export SSH_AGENT_PID;"+
			"echo Agent pid 4466;\n")

	agent := getAgetEnv(out)
	assert.Equal(t, lists.String("SSH_AUTH_SOCK=/tmp/ssh-kQN4kvgauzrv/agent.4465", "SSH_AGENT_PID=4466"), agent.env)
	assert.Equal(t, "/tmp/ssh-kQN4kvgauzrv/agent.4465", agent.socket)
}

func TestRunAgentAgent(t *testing.T) {
	os2.Setenv("SSH_AUTH_SOCK", "")
	os2.Setenv("SSH_AGENT_PID", "")
	err, agent := CheckSshAgentOrRun()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(agent.env))

	err = agent.AddSshKey(ssh_key, "")
	assert.Nil(t, err)

	err = agent.AddSshKey(ssh_with_pw, ssh_with_pw_pw)
	assert.Nil(t, err)

	err, lines := os.ExecCmd(os.ExecParams{Env: agent.env}, "ssh-add", "-L")
	assert.Nil(t, err)
	assert.Equal(t, lists.String(ssh_key_pub, ssh_with_pw_pub), lines)

	agent.Stop()
}

const ssh_key = "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAlwAAAAdzc2gtcn\nNhAAAAAwEAAQAAAIEA4Itqf68ss3tQoNjjpiy83/2Ohgcu+8fny07aWY3VYgnLl4dKiGZe\nwj8AhmTQ7URfXC8B5hndGtGC7rvyt/F1uDlVvdFyVzGZl/UUdaXF5yzx3a9hjmhlO1BqvQ\nUSeBEmZ99J5TRmHaoxwVQ8T1obS7cyMOHpmnUa0lwmxAIvSrcAAAIICkSBAQpEgQEAAAAH\nc3NoLXJzYQAAAIEA4Itqf68ss3tQoNjjpiy83/2Ohgcu+8fny07aWY3VYgnLl4dKiGZewj\n8AhmTQ7URfXC8B5hndGtGC7rvyt/F1uDlVvdFyVzGZl/UUdaXF5yzx3a9hjmhlO1BqvQUS\neBEmZ99J5TRmHaoxwVQ8T1obS7cyMOHpmnUa0lwmxAIvSrcAAAADAQABAAAAgAwXVcfEXg\nrYJBJVO4TyOcVx+N+8uUnzjMbE2zshSRE7Z8wkC95mbMnW7KdP/HQaT2w+V8LVN7O+/mbu\nlfZTuTv1fQfNhAPtzKfjsBaC6+Lf6dd6IcxD60N5o1yrD/tyyvtLYJ/HVqnmJMcWSx4eOB\nBmLyfIudO8u5XGkgOwPMa5AAAAQD2Flt3jx3avnQ2KD2gh0VvJgQGANCDmG7JNR3III8hJ\nKF0pg0VCcLun5zYdSKb8z8QV10h2D5JmfCj2jSdU9UwAAABBAPdXfwz8jFTG34FZwcOgSq\nQGFIL+CBUa0OBfegHmX+En8+Sqh2EGg8kjZxSq5VPP6pr/Gstfiuz/u+y0t2pmp8MAAABB\nAOhnoJYxECdH+HtGeev715kw48xSz6+GhjpyigABebMw3P5T+/gxTVyLsCH0xD3+S4wXYK\ntsW4PNblg/Ut4dlf0AAAARcm9vdEBkODNjOWEzNDViMzIBAg==\n-----END OPENSSH PRIVATE KEY-----\n"
const ssh_key_pub = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDgi2p/ryyze1Cg2OOmLLzf/Y6GBy77x+fLTtpZjdViCcuXh0qIZl7CPwCGZNDtRF9cLwHmGd0a0YLuu/K38XW4OVW90XJXMZmX9RR1pcXnLPHdr2GOaGU7UGq9BRJ4ESZn30nlNGYdqjHBVDxPWhtLtzIw4emadRrSXCbEAi9Ktw== root@d83c9a345b32"

const ssh_with_pw = "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAACmFlczI1Ni1jdHIAAAAGYmNyeXB0AAAAGAAAABBXtQ+AZK\nC6XpUqP+LYlataAAAAEAAAAAEAAAGXAAAAB3NzaC1yc2EAAAADAQABAAABgQDHqxsxt6Kz\nDr/MXkpiqGbb8NN7KIfR3/7lNUgrI9l2KozuI5Hr0F5e1NmjRctKsRRWi0SMnuuUGwGRSv\nuGEquW5uytkZdaKmWBtff8aTEw98VVLy1iOp92Bf4kbOI/3u2elohm0q9JkEoIx2Xt/W3P\nJgoqQQ9c839rREmDhlT3G6ows1qj5BlqArHR3dGRqCGq5d3W/coNlXUchglepGlvBh3yrx\n3BAu4ooGclVfgbjHlMz4U0TmhDvM9SEHkJ3OvjNbpqxXnRzEqDbrhvSctrErXnGc8/Ffpp\n8hK6TTuV9036YTTdzx4IldoaojYNLjhuQUfbfKcIk2rYfgsUqY1ymjiKjj59S8QCuM3uFZ\nkvDDlIwwrOOUlSR3LIgMCkhb9P/P2oeszS7Bw+BBwMaIowRKT9l2ZGWrV3b1Mur86TirJD\n0+oaK5b/m6uPegTb3p/UlJkT/g2/IaqH98tA41iitXaaiFk8szF0BOPl2chvsB9/0J23Wa\nnXkUQnTtVurKsAAAWQlhS+mdLIIiyiIbsc8v81gShAQzP5VxPG1ODPC/wyOLUqE/I2ZUNY\nmHZKhIUpi6UD0J9PyqqnvM04uzyvkLrdqpV/ZeJptFijLyxxS4SG3mehD4J8OWheTo0c3T\nzScZKflwafzlJg3BlsWs7U5viobIdNJKJd9hB5eriwf4mRkjwA3a+NWEUW+alGNSXyr/cI\n9/0jUeo1jaGxAQr0lhR6OinTs+umQ4NEaJLXQRxWnPoW6XKe4EdQLlsoSfBuYtcSE8b01W\ntvtSxFbGJGC+seXOv+uGTsIiTVSYj3XwJCCCv+dbo5viy38DHQO9kb1HVp7hjYNQNXkfTH\nzQlzJMWuhF5IEtK1nrrwfWFDGd+ps9jp65m5W/MhiX+fmpCghf5IOtHOvK3cDhD7QK3vBm\nxraXR3G692BngU+vBcodW4DuIW4A/qyr+HxTHYRr2WVoMJN4GgiH3HplnLS3/T9Nh7bqoL\n59wGqIbv6Dcv56EvD/Rc2Dza05HWK2hA4W3rc3W/u/6JodrGpKHyvcNpJbAq6yxrW9GY6G\n5M6/zWAzNxcVFh5Wx/YFDKnOeYUBZIIT8tvNlN0wiqqFyKhviek2OK6VO7859UVJV7LfNn\npW3vbhJ++gbOqrH/26DHegnU/Gs6ZOTIJ2d2SdlYQaC8+/euocSJSOKi5UGdCIS52rJ4lj\nDpcPZARagaDn4qrwne2UPHhZcINixKb6Tqnl+Y/bD6TTPl9vMUMfZDk6nSnySrZytl12oD\n2EMUHGRtiIqH+zvFZKjSqDbVdzhebWGsxo9C1r/6WGIP+S8Suze4HIjkAUaujyWfc7OO/A\n635GK13KSMA8SahxD+m8gmu4iRCFK5jYeOycr1DY/PU+IL8ZxyURJZjHiVeFXYN/e0pn7y\nN9mZMNk2k2OyMnIwHmDgHZZMBWvO5NK9xlcKUhc02XJcZVlj5Ks53sT2MroKAC/MbkCtiE\nGwULPx3mP2PMIDjYv3lIwRoG6+kpL6VnH44j4RH5vKhCwYkLyjx8Dt8ROZ9nMIX9GLoFLk\nY6yL1JhFIByY7SWgxMrO1XGKs9X4VQn2Rcwl+ZjgOK91xEp/1BZGRT4xKIAvMkH+QAqUaP\nhPxaFB5Lgt6OHufQX1DFxlhMaVgBOlOk+dbyRKgbT/8xG3ZoVnNb5MLYipw5+z5kqfWxNT\nzUxY61/8rJfE0mzgYmLELv9IRBEtR7G1aIckF9qkQ5c9iIUnB9Eis+eVD4fQuVt6vzUPDG\nrUc/E9u/ORazxQGAyRXvY8BIPrINswCPqF8ki6aJ+5oERm4M9uzMo5sLnnnI+ynH430Fzh\nXclv+Lct5NTRPcPp4Umn+yxHxvrXxj3ZYWrHk83kDtaESYIBCVLiGWTtY98Zc9YUdBRhN0\nRa524D/EY6Ne08J0YzI4shHey/h+Ajj7J2KYkN6Iy69AJgYE73g6SXKDrIJsrdY3QPto6y\nC3HSqZOclKyOyFRCKdrVrJZunjyKt3kxfYL4Dy3VDbLW8PvWBWrB3ZH5+HlXiw6MVVL+y4\n4ZofnthNfphDfo8PKBapKCoi6ORgFzyDPrg8XV8d1BS2PSRFMg2GdFC5nTD4mcjoNFqhki\n+jeoMa6rW/V5Nsa16oxH7flkQqrs9LE1EEUnYX1lgcUT98N7NBnfo4HeznPSn4l6HphbUQ\nF1Ngo5GzDiLkh3yKXrZfKyMzjIdG/t9t5BkAuAKmhC5jf4wyNOPpuAWj6GX58QvUpyZrrg\nnmjXYSsYFLhQUiQzaX4PpczAw+hbIKq8Zgcst3vlWvXAYK3PBgMn4v1oKwsdKNSFNbJtBB\n9oC13CrSQEeHAjSHqmjN6rtLXc0e/jML3Rf5GG7pTnaHlKZSanJOgIWChg3ZWf3p9D3771\nA2Ub19DE1aCDStzLVIRwnBbmhvA=\n-----END OPENSSH PRIVATE KEY-----\n"
const ssh_with_pw_pw = "1"
const ssh_with_pw_pub = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDHqxsxt6KzDr/MXkpiqGbb8NN7KIfR3/7lNUgrI9l2KozuI5Hr0F5e1NmjRctKsRRWi0SMnuuUGwGRSvuGEquW5uytkZdaKmWBtff8aTEw98VVLy1iOp92Bf4kbOI/3u2elohm0q9JkEoIx2Xt/W3PJgoqQQ9c839rREmDhlT3G6ows1qj5BlqArHR3dGRqCGq5d3W/coNlXUchglepGlvBh3yrx3BAu4ooGclVfgbjHlMz4U0TmhDvM9SEHkJ3OvjNbpqxXnRzEqDbrhvSctrErXnGc8/Ffpp8hK6TTuV9036YTTdzx4IldoaojYNLjhuQUfbfKcIk2rYfgsUqY1ymjiKjj59S8QCuM3uFZkvDDlIwwrOOUlSR3LIgMCkhb9P/P2oeszS7Bw+BBwMaIowRKT9l2ZGWrV3b1Mur86TirJD0+oaK5b/m6uPegTb3p/UlJkT/g2/IaqH98tA41iitXaaiFk8szF0BOPl2chvsB9/0J23WanXkUQnTtVurKs= root@d83c9a345b32"
