with r1_categories as (
	select category, air_date, round
	from jeopardy_clues
	group by category, air_date, round
	having count(*) = 5 and round = 1
	order by random() asc
	limit $1
),
r2_categories as (
	select category, air_date, round
	from jeopardy_clues
	group by category, air_date, round
	having count(*) = 5 and round = 2
	order by random() asc
	limit $2
),
round1 as (
	select jc.round, jc.clue_value, jc.category, jc.comments, jc.answer, jc.question, jc.alternatives, jc.air_date
	from jeopardy_clues as jc 
	join r1_categories as r1 
	on jc.category = r1.category and jc.air_date = r1.air_date and jc.round = r1.round
),
round2 as (
	select jc.round, jc.clue_value, jc.category, jc.comments, jc.answer, jc.question, jc.alternatives, jc.air_date
	from jeopardy_clues as jc 
	join r2_categories as r2
	on jc.category = r2.category and jc.air_date = r2.air_date and jc.round = r2.round
),
final_jeopardy as (
	select jc.round, jc.clue_value, jc.category, jc.comments, jc.answer, jc.question, jc.alternatives, jc.air_date
	from jeopardy_clues as jc
	where round = 3
	order by random()
	limit 1
),
all_clues as (
	select * from round1
	union all
	select * from round2
	union all
	select * from final_jeopardy
)
-- union all, not union: two clues in the same category can be identical rows, and
-- deduping them leaves the category short, which shifts every category after it.
-- air_date is in the sort so same-named categories from different episodes stay
-- contiguous instead of interleaving by clue_value.
select round, clue_value, category, comments, answer, question, alternatives
from all_clues
order by round asc, category asc, air_date asc, clue_value asc;
