// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package advance

import (
	"context"
	"fmt"

	"github.com/GoogleContainerTools/kpt/rollouts/api/v1alpha1"
	"github.com/GoogleContainerTools/kpt/rollouts/rolloutsclient"
	"github.com/spf13/cobra"
)

func newRunner(ctx context.Context) *runner {
	r := &runner{
		ctx: ctx,
	}
	c := &cobra.Command{
		Use:     "advance rollout-name wave-name",
		Short:   "advances the wave of a progressive rollout",
		Long:    "advances the wave of a progressive rollout",
		Example: "advances the wave of a progressive rollout",
		RunE:    r.runE,
	}
	r.Command = c
	return r
}

func NewCommand(ctx context.Context) *cobra.Command {
	return newRunner(ctx).Command
}

type runner struct {
	ctx     context.Context
	Command *cobra.Command
}

func (r *runner) runE(cmd *cobra.Command, args []string) error {
	rlc, err := rolloutsclient.New()
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("must provide rollout name\n")
	}

	if len(args) == 1 {
		return fmt.Errorf("must provide wave name\n")
	}

	rolloutName := args[0]
	waveName := args[1]

	rollout, err := rlc.Get(r.ctx, rolloutName)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	if rollout.Spec.Strategy.Type != v1alpha1.Progressive {
		return fmt.Errorf("rollout must be using the progressive strategy to use this command\n")
	}

	if rollout.Status.WaveStatuses != nil {
		waveFound := false

		for _, waveStatus := range rollout.Status.WaveStatuses {
			if waveStatus.Name == waveName {
				waveFound = true
				break
			}
		}

		if !waveFound {
			return fmt.Errorf("wave %q not found in this rollout\n", waveName)
		}
	}

	rollout.Spec.Strategy.Progressive.PauseAfterWave.WaveName = waveName

	err = rlc.Update(r.ctx, rollout)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	fmt.Println("done")
	return nil
}
