# Upkeep

A console application for time tracking.

Install locally with `make install` or build an executable with `make build-prod`.

#### Data storage

Configuration and past blocks are stored as JSON files in the `.upkeep` folder in the user home directory.

## Usage

Upkeep tracks time in blocks on a per-day basis. Each block belongs to a category. You can view totals per day or category and create exports of time ranges.

Run `upkeep help` to view the available commands.

### Examples

- Every day when you start working on some task, you start upkeep, specifying the task category.

```upkeep start mytask```

- When you start another task, just start a new task.

```upkeep start othertask```

- When you want to temporarily switch to a different task but come back to the current task, you can switch to that task.

```upkeep switch smallthing```

- And afterwards continue with the previous task.

```upkeep continue```

- When you're done for that day, you end the final block.

```upkeep stop```

- If you want to add a block that has no specific start or stop time but should add to a days total, write it specifically.

```upkeep write holiday 8h```

- Finalise this weeks timesheets, preventing them from further editing. See below for the date parameter (`d:w`) syntax.

```upkeep finalise d:w```

- Write an export of the current month.

```upkeep export d:m```

- (Optional) You determine how many hours you work each workday and set them as a quotum. For example on monday (day 1) you expect to work six hours. This only has to be done once.

```upkeep conf quotum 1 6h```

- Adjust (or remove with `none`) a (pre defined) quotum for a single day.

```upkeep quotum 3h```

- (Optional) If you want a specific category to only add to the day total for a certain maximum of time, set a category quotum.

```upkeep cat quotum break 30m```

### Date parameter

Most commands accept a date (range) parameter given as `d:<date>`.
Here `<date>` can be one of the following:
- Any date in the form of `YYYY-MM-DD` to select that specific date.
- `d` or `day` to select a day (default). Without shifting, it selects the current day.
- `w` or `week` to select a week. Without shifting, it selects the current week.
- `wr` or `weekrolling` to select a week by rolling window, ending on the current day.
- `m` or `month` to select a month. Without shifting, it selects the current month.

The letter/word inputs can be prefixed with (negative) numbers to shift their time range. For instance to select two weeks ago, use `d:-2w`.

### Group parameter

Categories for a block can be specified in two parts with a point separating the parts. In this case, the first part is treated as the group name. This allows for grouping block times in outputs/exports.

```upkeep start group.task```

If you want to view the grouped times, specify the group parameter `g:<group>`. Here `<group>` can be one of the following:
- `g` to use groups
- `c` to use the categories (default).

```upkeep view cat g:g```

## Development

Build development executable with `make build` or a debug executable with `make build-dbg`.

The development version will use the `./data` folder to read and write its data.