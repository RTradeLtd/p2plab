package experiments

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	parser "github.com/Netflix/p2plab/cue/parser"
	"github.com/Netflix/p2plab/metadata"
	"github.com/google/uuid"
)

func TestExperimentDefinition(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, cleanup := newTestDB(t, "exptestdir")
	defer func() {
		if err := cleanup(); err != nil {
			t.Fatal(err)
		}
	}()
	var ids []string
	t.Run("Experiment Creation And Retrieval", func(t *testing.T) {
		sourceFiles := []string{
			"../cue/cue.mod/p2plab_example1.cue",
			"../cue/cue.mod/p2plab_example2.cue",
		}
		for _, sourceFile := range sourceFiles {
			name := strings.Split(sourceFile, "/")
			exp1 := newTestExperiment(t, sourceFile, name[len(name)-1])
			ids = append(ids, exp1.ID)
			exp2, err := db.CreateExperiment(ctx, exp1)
			if err != nil {
				t.Fatal(err)
			}
			if exp1.ID != exp2.ID {
				t.Fatal("bad id")
			}
			if exp1.Status != exp2.Status {
				t.Fatal("bad status")
			}
			if !reflect.DeepEqual(exp1.Definition, exp2.Definition) {
				t.Fatal("bad definition")
			}
			exp3, err := db.GetExperiment(ctx, exp1.ID)
			if err != nil {
				t.Fatal(err)
			}
			if exp1.ID != exp3.ID {
				t.Fatal("bad id")
			}
			if exp1.Status != exp3.Status {
				t.Fatal("bad status")
			}
			if len(exp1.Definition.TrialDefinition) != len(exp3.Definition.TrialDefinition) {
				fmt.Println(exp1.Definition.TrialDefinition)
				fmt.Println(exp3.Definition.TrialDefinition)
				t.Fatal("bad trial definitions returned")
			}
		}
	})
	t.Run("List Experiments", func(t *testing.T) {
		experiments, err := db.ListExperiments(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for _, experiment := range experiments {
			if experiment.ID != ids[0] && experiment.ID != ids[1] {
				t.Fatal("bad experiment id found")
			}
		}
	})
}

func newTestExperiment(t *testing.T, sourceFile, name string) metadata.Experiment {
	data, err := ioutil.ReadFile("../cue/cue.mod/p2plab.cue")
	if err != nil {
		t.Fatal(err)
	}
	sourceData, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		t.Fatal(err)
	}
	psr := parser.NewParser([]string{string(data)})
	inst, err := psr.Compile(
		name,
		string(sourceData),
	)
	if err != nil {
		t.Fatal(err)
	}
	edef, err := inst.ToExperimentDefinition()
	if err != nil {
		t.Fatal(err)
	}
	return metadata.Experiment{
		ID:         uuid.New().String(),
		Status:     metadata.ExperimentRunning,
		Definition: edef,
	}
}

func newTestDB(t *testing.T, path string) (metadata.DB, func() error) {
	db, err := metadata.NewDB(path)
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() error {
		if err := db.Close(); err != nil {
			return err
		}
		return os.RemoveAll(path)
	}
	return db, cleanup
}
