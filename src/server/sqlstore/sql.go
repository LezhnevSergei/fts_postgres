package sqlstore

const SaveTSVector = `
insert into incident_tvectors (incident_id, tvs)
select incident_id, make_tsvector(r.display_name, r.description)
from incidents
         left join rules r on incidents.rule_id = r.rule_id;
`

const ListRules = `
select i.incident_id,
       r.display_name,
       r.description
from roswell.incidents as i
left join roswell.rules r on i.rule_id = r.rule_id
limit 50;
`

const SearchIncidents = `
select i.incident_id
from plainto_tsquery($1) as q,
     roswell.incidents as i
         left join roswell.tsvectors as r_ts on i.incident_id = r_ts.incident_id
where r_ts.tvs @@ q
order by ts_rank(r_ts.tvs, q) desc
limit 50;
`

const CalcSearchIncidents = `
explain (ANALYZE, FORMAT JSON)
select inc.incident_id
from plainto_tsquery($1) as q, roswell.incidents as inc
         left join roswell.tsvectors as r_ts on inc.incident_id = r_ts.incident_id
where r_ts.tvs @@ q
order by ts_rank(r_ts.tvs, q) desc
limit 50;
`
