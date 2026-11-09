-- A category is only usable when its five clues carry the exact slot values for the
-- round, and its round has all 30 clues. The scraper assigned categories by counting
-- clues, so a missing clue shifted every later clue in that round into the wrong
-- category. Those categories can still look complete, so the whole round is skipped.
with full_rounds as (
	select air_date, round
	from jeopardy_clues
	where round in (1, 2)
	group by air_date, round
	having count(*) = 30
),
valid_categories as (
	select jc.category, jc.air_date, jc.round
	from jeopardy_clues as jc
	join full_rounds as fr on jc.air_date = fr.air_date and jc.round = fr.round
	group by jc.category, jc.air_date, jc.round
	having array_agg(jc.clue_value order by jc.clue_value) =
		case when jc.round = 1
			then array[200, 400, 600, 800, 1000]
			else array[400, 800, 1200, 1600, 2000] end
),
r1_categories as (
	select category, air_date, round
	from valid_categories
	where round = 1
	order by random() asc
	limit $1
),
r2_categories as (
	select category, air_date, round
	from valid_categories
	where round = 2
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
