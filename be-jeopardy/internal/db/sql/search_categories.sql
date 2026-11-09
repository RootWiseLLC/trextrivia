-- Same validity rule as get_questions.sql: exact slot values, and a round with all 30 clues.
with full_rounds as (
	select air_date, round
	from jeopardy_clues
	where round in (1, 2)
	group by air_date, round
	having count(*) = 30
)
select jc.category, jc.round, jc.air_date
from jeopardy_clues as jc
join full_rounds as fr on jc.air_date = fr.air_date and jc.round = fr.round
where lower(jc.category) like concat($3::text, $1::text, '%')
group by jc.category, jc.round, jc.air_date
having (jc.round = 1 or jc.round = $2)
	and array_agg(jc.clue_value order by jc.clue_value) =
		case when jc.round = 1
			then array[200, 400, 600, 800, 1000]
			else array[400, 800, 1200, 1600, 2000] end
order by jc.category asc;
