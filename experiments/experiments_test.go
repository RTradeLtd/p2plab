package experiments

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
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
	exp1 := newTestExperiment(t)
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
	if !reflect.DeepEqual(exp1.Definition, exp3.Definition) {
		t.Fatal("bad definition")
	}
}

func newTestExperiment(t *testing.T) metadata.Experiment {
	data, err := ioutil.ReadFile("../cue/cue.mod/p2plab.cue")
	if err != nil {
		t.Fatal(err)
	}
	sourceData, err := ioutil.ReadFile("../cue/cue.mod/p2plab_example1.cue")
	if err != nil {
		t.Fatal(err)
	}
	psr := parser.NewParser([]string{string(data)})
	inst, err := psr.Compile(
		"expdeftest",
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
