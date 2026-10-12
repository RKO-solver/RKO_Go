# `tuning` Package

Chooses a metaheuristic's parameters while the search runs, learning from how well each choice worked.

The best population size, elite ratio or mutation rate depends on the instance and on the stage of the search. Instead of fixing them up front, a metaheuristic splits its run into cycles (a generation, a batch of iterations, a restart, ...). Before each cycle it asks a `Tuner` which configuration to use; after the cycle it reports a reward, and the tuner uses the reward to choose better next time.

## Quick start

```go
package main

import (
	"github.com/RKO-solver/rko-go/metaheuristic/tuning"
	"github.com/RKO-solver/rko-go/random"
)

func main() {
	space := tuning.NewSpace()
	popSize := tuning.Add(space, []int{100, 500, 1000})
	elite := tuning.Add(space, []float64{0.10, 0.20})
	mutant := tuning.Add(space, []float64{0.05, 0.10, 0.15})

	tuner := tuning.NewQLearning(space, random.New(42))

	const cycles = 1000
	for cycle := range cycles {
		state := tuner.Choose(float64(cycle) / cycles)

		// runCycle is your code: one generation with these values, returning its reward.
		reward := runCycle(popSize.At(state), elite.At(state), mutant.At(state))

		tuner.Learn(reward, float64(cycle+1)/cycles)
	}
}
```

1. Register every parameter with `Add`, then create the tuner. The tuner sizes its tables when it is created and does not see parameters added later.
2. Each cycle, call `Choose`, read the values with `At`, run the cycle, then call `Learn`.
3. Always call `Choose` before `Learn`: `Learn` rewards the configuration returned by the most recent `Choose`. A `Learn` before the first `Choose` updates the wrong estimate, and nothing detects it.

## Types

| Name | Is |
|---|---|
| `Space` | The set of configurations: every combination of the parameters' values. Each configuration is an `int`, its **state**, in `[0, States())`. |
| `NewSpace()` | An empty `Space` (one state, no parameters). The zero `Space` is not usable. |
| `Add(space, values)` | Registers a parameter and returns its `Param[T]`. `T` is any comparable type (`int`, `float64`, `string`, ...). |
| `Param[T].At(state)` | The parameter's value in `state`, already of type `T`. `state` is not checked. |
| `Tuner` | The interface a metaheuristic depends on: `Choose(progress) int` and `Learn(reward, progress)`. |
| `QLearning` | The `Tuner` implemented here: tabular Q-learning with epsilon-greedy exploration. |
| `NewQLearning(space, generator, options...)` | Creates a `QLearning`. |
| `WithDiscount(d)` | Option that sets the discount (default `DefaultDiscount`, `0.8`). |
| `Reporter` | Optional interface for tuners that can describe what they learned: `Policy() logger.Table`. A `Tuner` does not have to implement it. |
| `QLearning.Policy()` | Implements `Reporter`: the best move from each state (see [Reading the learned policy](#reading-the-learned-policy)). |

Write metaheuristics against `Tuner` and `Param`, not `QLearning`, so another learning method can replace it without changing them.

A `Param` is a small value that can be copied freely, and it stays valid when more parameters are added to its `Space`. `Add` keeps the `values` slice without copying it: do not modify it afterwards.

## Progress and reward

**Progress** is the fraction of the run already done, in `[0, 1]`: `0` at the start, `1` at the end. Compute it however suits the metaheuristic, for example elapsed time over the time limit, or iterations over the iteration limit. Pass the progress before the cycle to `Choose` and after it to `Learn`. Tuners trust it: it is neither clamped nor checked. `QLearning` uses it to explore less and learn more slowly as the run advances.

**Reward** says how well the cycle went; higher is better. The metaheuristic defines it. For example: `1` when the cycle found a new best solution, otherwise the relative gap to the best solution, which is at most `0`.

## How states are numbered

A state packs one value index per parameter into a single number. It works like a car's odometer whose wheels have different numbers of positions: the first parameter added is the fastest wheel, and each later wheel advances one position every time all the wheels before it complete a full turn.

The **stride** of a parameter is how much the state changes when that parameter moves to its next value: the product of the value counts of the parameters added before it. Then

```text
state = index0*stride0 + index1*stride1 + index2*stride2 + ...
```

With the three parameters of the quick start:

| Parameter | Values | Count | Stride |
|---|---|---|---|
| `popSize` | `100, 500, 1000` | 3 | 1 |
| `elite` | `0.10, 0.20` | 2 | 3 |
| `mutant` | `0.05, 0.10, 0.15` | 3 | 6 |

There are `3 × 2 × 3 = 18` states, and `state = pop + 3*elite + 6*mutant`, where `pop`, `elite` and `mutant` are value indices:

| State | `popSize` | `elite` | `mutant` |
|---|---|---|---|
| 0 | 100 | 0.10 | 0.05 |
| 1 | 500 | 0.10 | 0.05 |
| 2 | 1000 | 0.10 | 0.05 |
| 3 | 100 | 0.20 | 0.05 |
| ... | | | |
| 10 | 500 | 0.20 | 0.10 |
| ... | | | |
| 17 | 1000 | 0.20 | 0.15 |

`At` reverses this: dividing by the stride drops the faster wheels, and the remainder by the count drops the slower ones. For state 10:

```text
popSize = values[10/1 % 3] = values[1] = 500
elite   = values[10/3 % 2] = values[1] = 0.20
mutant  = values[10/6 % 3] = values[1] = 0.10
```

## How `QLearning` works

`QLearning` keeps an estimate `Q(s, a)` of how good it is to take action `a` in state `s`. Each `Choose` takes one action from the current state and returns the state it reaches, the configuration for the next cycle. The following `Learn` updates the estimate of that move from the reward and makes the reached state current.

### Actions

An action either keeps the configuration (**stay**) or changes exactly one parameter to another of its values, so the tuner moves through the space one parameter at a time. Every state has `1 + Σ(count − 1)` actions, laid out as:

- action 0: stay;
- then one block per parameter, in the order they were added, listing that parameter's other values in increasing index order and skipping the current one.

With the example above there are `1 + 2 + 1 + 2 = 6` actions. From state 10 = `(500, 0.20, 0.10)`:

| Action | Change | Next state |
|---|---|---|
| 0 | stay | 10 |
| 1 | `popSize` → 100 | 9 |
| 2 | `popSize` → 1000 | 11 |
| 3 | `elite` → 0.10 | 7 |
| 4 | `mutant` → 0.05 | 4 |
| 5 | `mutant` → 0.15 | 16 |

Changing one parameter changes one digit of the state, so the next state is `state + (newIndex − currentIndex) × stride`, just as turning 345 into 375 adds `(7 − 4) × 10`. For action 3: `10 + (0 − 1) × 3 = 7`. All these transitions are computed once in `NewQLearning`, so `Choose` moves with a table lookup.

The Q-table holds `States() × actions` values, so memory and creation time grow with the product of the parameters' value counts.

### Choosing and learning

`Choose(progress)` is epsilon-greedy: with probability `epsilon(progress)` it takes a uniformly random action (which may also be the best one or stay), otherwise the action with the highest `Q`. Ties go to the lowest action, so with equal estimates the tuner stays where it is. Calling `Choose` twice without `Learn` replaces the first choice.

`Learn(reward, progress)` applies the Q-learning update to the last choice:

```text
Q(s, a) += alpha * (reward + gamma * max Q(s', ·) − Q(s, a))
```

`s` is the current state, `a` the chosen action, `s'` the state `Choose` returned, `alpha` the learning rate and `gamma` the discount. The target is the reward plus the discounted value of the best move from `s'`; the update moves `Q(s, a)` the fraction `alpha` of the way towards it. Then `s'` becomes the current state.

### Exploration rate

The exploration rate `epsilon` follows a cosine schedule with warm restarts:

- Progress is split into slices of `0.1`.
- Within a slice, `epsilon` eases from the slice's peak down to `0.1` along half a cosine: slowly at both ends, fastest in the middle.
- At the start of the next slice it jumps back to a peak `0.1` lower: `1.0`, `0.9`, `0.8`, ..., down to `0.1`. From progress `0.9` on it stays at `0.1`.

So the tuner explores widely early on, and each restart explores a little less.

| Progress | Slice | Peak | `epsilon` |
|---|---|---|---|
| 0 | 0 | 1.0 | 1.0 |
| 0.05 | 0 | 1.0 | 0.55 |
| 0.15 | 1 | 0.9 | 0.5 |
| 0.35 | 3 | 0.7 | 0.4 |
| 1 | 10 | 0.1 | 0.1 |

`epsilon` depends on progress alone, so it is correct even when one cycle spans several slices.

### Learning rate

`alpha = 1 − 0.9 × progress`: it falls linearly from `1.0` at the start, where an update replaces `Q(s, a)` by its target, to `0.1` at the end, where each reward only nudges the estimates.

### Discount

The discount `gamma` weights the value of the state an action leads to against the immediate reward, so the tuner also favours configurations that lead to good next configurations. The default is `DefaultDiscount` (`0.8`); change it with an option:

```go
tuner := tuning.NewQLearning(space, rg, tuning.WithDiscount(0.5))
```

With `0` the tuner learns only from immediate rewards; values closer to `1` weight the states reached more. The discount must be in `[0, 1)`: at `1` or above, Q-values can grow without bound, since the run has no fixed end.

### Initial values and reproducibility

Initial Q-values are drawn uniformly from `[0.005, 0.01)`. Cycles that do not improve typically earn a reward of at most `0`, so this small positive start makes untried actions look slightly better than tried ones that did not improve, and gets them tried. The random values also make each state's first best action random. The tuner starts in a uniformly random state.

All randomness comes from the generator passed to `NewQLearning`, so the same seed gives the same choices for the same rewards and progress. Options do not draw random numbers.

## Reading the learned policy

A tuner that implements `Reporter` returns what it learned as a `logger.Table`. For `QLearning`, `Policy()` has one row per state, with these columns:

| Column | Holds |
|---|---|
| `state` | The state number. |
| `p0`, `p1`, ... | Each parameter's value in that state, in the order the parameters were added, printed with `fmt.Sprint` (so `0.10` shows as `0.1`). Parameters have no names, only these numbers. |
| `move` | The action with the highest `Q` from that state: `stay`, or `p<d> -> <value>` for the one parameter it changes. |
| `next` | The state that move reaches. |
| `q` | That action's Q-value, as the shortest decimal that reads back exactly. |

The table describes the `Space` as it was when the tuner was created; parameters added later are not in it. Part of the policy learned for the quick-start space:

```text
state  p0    p1   p2    move        next  q
0      100   0.1  0.05  p1 -> 0.2   3     2.8730090733630296
1      500   0.1  0.05  p1 -> 0.2   4     3.552126594947072
...
9      100   0.2  0.1   p0 -> 1000  11    3.8999581614054986
10     500   0.2  0.1   stay        10    4.999999999999994
11     1000  0.2  0.1   p0 -> 500   10    4.999999999999994
...
```

Following `next` from any row shows where the tuner would go if it stopped exploring; a `stay` row is a configuration it would keep.

To print or save it, pass it to the solver's logger at the end of the run; check for `Reporter`, since a `Tuner` may not implement it:

```go
if r, ok := tuner.(tuning.Reporter); ok {
	solverLogger.Policy(r.Policy())
}
```

Then, after `Solve`, call `log.PrintPolicies(os.Stdout)` or `log.SavePoliciesCSV("")` on the logger. See [Policies](../../logger/README.md#policies) in the `logger` package.

## Concurrency

A `QLearning` is not safe for concurrent use: create one per solver goroutine, each with its own `random.Generator` (for example from `rg.Split()`). Once every parameter is added, a `Space` and its `Param`s are only read, so they can be shared.

## Panics

| Call | Panics when |
|---|---|
| `Add` | `values` is empty, or holds the same value twice (compared with `==`), since two states would then share one configuration. |
| `WithDiscount` | The discount is outside `[0, 1)` or NaN. It panics when called, not when `NewQLearning` applies it. |

## See also

- [Project README](../../README.md)
- [`metaheuristic` package](../README.md)
- [`logger` package](../../logger/README.md)
