package role

func TestDelegateVotes_writedb(t *testing.T) {
	ins := db.NewDbService("./file2", "./file2/db.db")
	err := ApplyPersistanceRole(ins, block)
	if err != nil {
		fmt.Println(err)
	}
	value := &DelegateVotes{
		OwnerAccount: "nodepad",
		Serve: Serve{
			Votes:          1,
			Position:       big.NewInt(2),
			TermUpdateTime: big.NewInt(2),
			TermFinishTime: big.NewInt(2),
		},
	}
	err = SetDelegateVotesRole(ins, value.OwnerAccount, value)
	if err != nil {
		fmt.Println("SetDelegateVotesRole", err)
	}

	value, err = GetDelegateVotesRoleByAccountName(ins, value.OwnerAccount)
	if err != nil {
		fmt.Println("GetDelegateVotesRoleByAccountName", err)
	}
	fmt.Println(value)

	value, err = GetDelegateVotesRoleByVote(ins, value.Serve.Votes)
	if err != nil {
		fmt.Println("GetDelegateVotesRoleByVote", err)
	}
	fmt.Println(value)

	value, err = GetDelegateVotesRoleByFinishTime(ins, value.Serve.TermFinishTime)
	if err != nil {
		fmt.Println("GetDelegateVotesRoleByFinishTime", err)
	}
	fmt.Println(value)

	values, nerr := GetAllDelegateVotes(ins)
	if nerr != nil {
		fmt.Println("GetAllDelegateVotes", nerr)
	}
	fmt.Println(len(values))

	svotes, nerr := GetAllSortVotesDelegates(ins)
	if nerr != nil {
		fmt.Println("GetAllSortVotesDelegates", nerr)
	}
	fmt.Println(len(svotes))
	tvotes, nerr := GetAllSortFinishTimeDelegates(ins)
	if nerr != nil {
		fmt.Println("GetAllSortFinishTimeDelegates", nerr)
	}
	fmt.Println(len(tvotes))
}
