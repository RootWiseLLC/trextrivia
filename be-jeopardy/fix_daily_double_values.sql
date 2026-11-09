-- The scraper stored a Daily Double's wager as its clue_value, so a $600 slot could
-- read $3,000. That shows the wrong value on the board and, because clues are sorted
-- by clue_value, moves the clue to the wrong row. For every category with five clues
-- and one Daily Double, set the Daily Double's value to the one slot value that the
-- other four clues leave empty. Safe to run more than once.
with cats as (
	select air_date, round, category,
		array_agg(clue_value) filter (where daily_double_value = 0) as vals
	from jeopardy_clues
	where round in (1, 2)
	group by air_date, round, category
	having count(*) = 5 and count(*) filter (where daily_double_value > 0) = 1
),
missing as (
	select c.air_date, c.round, c.category, min(v) as value
	from cats as c,
	unnest(case when c.round = 1
		then array[200, 400, 600, 800, 1000]
		else array[400, 800, 1200, 1600, 2000] end) as v
	where not (v = any(c.vals))
	group by c.air_date, c.round, c.category
	having count(*) = 1
)
update jeopardy_clues as jc
set clue_value = m.value
from missing as m
where jc.air_date = m.air_date
	and jc.round = m.round
	and jc.category = m.category
	and jc.daily_double_value > 0;
