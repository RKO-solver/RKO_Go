package brkga

import (
	"context"
	"log/slog"
	"math"
	"slices"
	"time"

	"github.com/RKO-solver/rko-go/definition"
	"github.com/RKO-solver/rko-go/metaheuristic"
	"github.com/RKO-solver/rko-go/metaheuristic/internal/genetic"
	"github.com/RKO-solver/rko-go/metaheuristic/rk"
	"github.com/RKO-solver/rko-go/metaheuristic/solution"
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
)

func (brkga *BRKGA) solve(ctx context.Context, solutionPool *solution.Pool) (*metaheuristic.RandomKeyValue, time.Duration) {
	ctx, ctxAction := definition.NewSolverContext(ctx)
	defer ctxAction.Stop()

	configuration := brkga.configuration
	env := brkga.env
	rg := brkga.rg
	local := brkga.search
	numKeys := env.NumKeys()
	tuner := brkga.tuner
	populationSize := configuration.PopulationSize[0]
	eliteRatio := configuration.EliteRatio[0]
	mutantRatio := configuration.MutantRatio[0]
	crossoverAlpha := configuration.CrossoverAlpha[0]
	mutationAlpha := configuration.MutationAlpha[0]
	maxGenerationNoImprovement := configuration.MaxGenerationNoImprovement[0]
	maxSize := slices.Max(configuration.PopulationSize)

	population := make([]*metaheuristic.RandomKeyValue, maxSize)
	populationIntermediary := make([]*metaheuristic.RandomKeyValue, maxSize)
	for i := range maxSize {
		population[i] = metaheuristic.NewRandomKeyValue(numKeys)
		populationIntermediary[i] = metaheuristic.NewRandomKeyValue(numKeys)
	}

	generationNoImprovement := 0

	start := definition.StartTime(ctx)
	bestPerson := genetic.InitialPopulation(population[:populationSize], env, true, rg)
	solutionPool.Add(bestPerson, start)

	deadline, ok := ctx.Deadline()
	total := deadline.Sub(start)
	ok = ok && total > 0
	n := ctxAction.Limit()
	progress := func() float64 {
		if ok {
			return min(1, float64(time.Since(start))/float64(total))
		}
		if n > 0 {
			return min(1, float64(ctxAction.Iteration())/float64(n))
		}
		l := math.Log1p(float64(ctxAction.Iteration()) / float64(brkga.states))
		return l / (1 + l)
	}
	if tuner != nil {
		defer func() {
			if r, ok := tuner.(tuning.Reporter); ok {
				brkga.logger.Policy(r.Policy())
			}
		}()
	}

	eliteSize := int(eliteRatio * float64(populationSize))
	mutantSize := int(mutantRatio * float64(populationSize))
	for ; ctx.Err() == nil; ctxAction.Step() {
		if tuner != nil {
			state := tuner.Choose(progress())
			previousSize := populationSize
			populationSize = brkga.populationSize.At(state)
			eliteRatio = brkga.eliteRatio.At(state)
			mutantRatio = brkga.mutantRatio.At(state)
			crossoverAlpha = brkga.crossoverAlpha.At(state)
			mutationAlpha = brkga.mutationAlpha.At(state)
			maxGenerationNoImprovement = brkga.maxGenerationNoImprovement.At(state)
			for i := previousSize; i < populationSize; i++ {
				rk.Reset(population[i].RK, rg)
				population[i].Cost = env.Cost(population[i].RK)
			}
			if populationSize > previousSize {
				metaheuristic.Sort(population[:populationSize])
			}
			eliteSize = int(eliteRatio * float64(populationSize))
			mutantSize = int(mutantRatio * float64(populationSize))
		}

		for i := range eliteSize {
			rk.Copy(populationIntermediary[i], population[i])
		}

		for i := eliteSize; i < (populationSize - mutantSize); i++ {

			// crossoverUniformElite
			eliteParentId := rg.IntN(eliteSize)
			nonEliteParentId := rg.RangeInt(eliteSize, populationSize)

			copy(populationIntermediary[i].RK, population[eliteParentId].RK)

			for j := range numKeys {
				if rg.Float64() < mutationAlpha {
					populationIntermediary[i].RK[j] = rg.Float64()
				} else {
					if rg.Float64() >= crossoverAlpha {
						populationIntermediary[i].RK[j] = population[nonEliteParentId].RK[j]
					}
				}
			}

			populationIntermediary[i].Cost = env.Cost(populationIntermediary[i].RK)
			// end crossoverUniformElite

		}

		for i := populationSize - mutantSize; i < populationSize; i++ {
			rk.Reset(populationIntermediary[i].RK, rg)
			populationIntermediary[i].Cost = env.Cost(populationIntermediary[i].RK)
		}

		population, populationIntermediary = populationIntermediary, population

		k := rg.IntN(populationSize)
		local.Search(ctx, population[k])

		metaheuristic.Sort(population[:populationSize])

		if ctx.Err() != nil {
			if population[0].Cost < bestPerson.Cost {
				rk.Copy(bestPerson, population[0])
			}

			solutionPool.Add(bestPerson, start)

			break
		}

		elapsed := time.Since(start)
		if population[0].Cost < bestPerson.Cost {
			rk.Copy(bestPerson, population[0])
			generationNoImprovement = 0
		} else {
			generationNoImprovement++
			brkga.logger.Event(slog.Attr{Key: "stalled", Value: slog.GroupValue(slog.Int("generation", ctxAction.Iteration()), slog.Int("no-improvement", generationNoImprovement))}, elapsed)

		}
		solutionPool.Add(bestPerson, start)
		brkga.logger.Register(population[0].Cost, bestPerson.Cost, ctxAction.Iteration(), elapsed)

		if tuner != nil {
			reward := 1 + 1/float64(populationSize)
			if generationNoImprovement > 0 {
				reward = float64(bestPerson.Cost-population[0].Cost) / math.Max(1, math.Abs(float64(bestPerson.Cost)))
			}
			tuner.Learn(reward, progress())
		}

		if generationNoImprovement > maxGenerationNoImprovement {
			brkga.logger.Event(slog.String("restart", "limit generations without improvement"), time.Since(start))
			restartBest := genetic.InitialPopulation(population[:populationSize], env, true, rg)
			if restartBest.Cost < bestPerson.Cost {
				rk.Copy(bestPerson, restartBest)
			}
			solutionPool.Add(bestPerson, start)

			generationNoImprovement = 0
		}

	}

	return bestPerson, time.Since(start)
}
