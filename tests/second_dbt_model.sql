
-- Use the `ref` function to select from other models

select *
from {{ ref('my_first_dbt_model') }}
join {{ ref('my_third_dbt_model') }}
where id = 1
