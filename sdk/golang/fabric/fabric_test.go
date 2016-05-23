package fabric

import (
	"testing"
	"fmt"
	"os"
	"time"
)

var (
	chain Chain
	alice Member
	bob Member
	charlie Member

	bobCredential Credential
	charlieCredential Credential
)

func TestMain(m *testing.M) {
	// Get Chain
	chain, err := GetChain(".")
	if err != nil {
		clientSDKLog.Error("Failed getting chain [%s]\n", err)
		panic(fmt.Errorf("Failed getting chain [%s].", err))
	}

	// Get members
	alice = chain.EnrollMember("alice", "jim", "6avZQLwcUe9b")
	bob = chain.EnrollMember("bob", "lukas", "NPKYL39uKbkj")
	charlie = chain.EnrollMember("charlie", "diego", "DRJ23pEQl16a")

	if err := chain.Flush(); err != nil {
		fmt.Printf("Failed initializing clients [%s]\n", err)
		panic(fmt.Errorf("Failed initializing clients [%s].", err))
	}

	// Prepare credentials
	bobCredential = bob.GetChaincodeByAlias("assetmgm").BindCredential()
	charlieCredential = charlie.GetChaincodeByAlias("assetmgm").BindCredential()

	// Run tests
	os.Exit(m.Run())
}

func TestAssetManagementDeploy(t *testing.T) {
	tx := alice.Deploy("assetmgm", "assetmgm")
	tx.WithCredential()
	tx.Confidential()
	tx.Send()
	if err := alice.Flush(); err != nil {
		t.Fatalf("Failed invoking [%s]", err)
	}
	if err := tx.Flush(); err != nil {
		t.Fatalf("Failed assigning [%s]", err)
	}
}

func TestAssetManagementAssign(t *testing.T) {
	time.Sleep(3 * time.Second)
	tx := alice.Invoke("assetmgm", "assign")
	tx.WithCredential()
	tx.AddArgument("Ferrari")
	tx.AddArgumentCredential(bobCredential)
	tx.Confidential()
	tx.Send()
	if err := alice.Flush(); err != nil {
		t.Fatalf("Failed invoking [%s]", err)
	}
	if err := tx.Flush(); err != nil {
		t.Fatalf("Failed assigning [%s]", err)
	}
}

func TestAssetManagementTransfer(t *testing.T) {
	time.Sleep(3 * time.Second)
	tx := bob.Invoke("assetmgm", "transfer")
	tx.WithCredential()
	tx.AddArgument("Ferrari")
	tx.AddArgumentCredential(charlieCredential)
	tx.Confidential()
	tx.Send()
	if err := bob.Flush(); err != nil {
		t.Fatalf("Failed transferring [%s]", err)
	}
	if err := tx.Flush(); err != nil {
		t.Fatalf("Failed transferring [%s]", err)
	}
}

func TestAssetManagementQuery(t *testing.T) {
	time.Sleep(3 * time.Second)
	tx := alice.Query("assetmgm", "query")
	tx.WithCredential()
	tx.AddArgument("Ferrari")
	tx.Confidential()
	tx.Send()
	result := tx.GetResponse()
	if err := alice.Flush(); err != nil {
		t.Fatalf("Failed transferring [%s]", err)
	}
	if err := tx.Flush(); err != nil {
		t.Fatalf("Failed transferring [%s]", err)
	}

	fmt.Printf("Result [% x]", result)
}