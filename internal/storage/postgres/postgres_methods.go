package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bondarenkozahar80-hub/3.6/internal/model"
)

var sortFields = map[string]bool{"id": true, "name": true, "amount": true, "type": true, "category": true, "event_date": true}

func (p *Postgres) Create(ctx context.Context, t *model.Transaction) (int, error) {
	query := `
        INSERT INTO transactions (name, description, amount, type, category, event_date)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id;
    `
	var id int
	err := p.DB.QueryRowContext(ctx, query, t.Name, t.Description, t.Amount, t.Type, t.Category, t.EventDate).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *Postgres) GetByID(ctx context.Context, id int) (*model.Transaction, error) {
	query := `
        SELECT  * FROM transactions
        WHERE id = $1;
    `
	var t model.Transaction
	err := p.DB.GetContext(ctx, &t, query, id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (p *Postgres) GetAll(ctx context.Context) ([]model.Transaction, error) {
	query := `
        SELECT * FROM transactions
        ORDER BY event_date DESC;
    `
	var list []model.Transaction
	err := p.DB.SelectContext(ctx, &list, query)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *Postgres) Update(ctx context.Context, t *model.Transaction) error {
	query := `
        UPDATE transactions
        SET name = $1,
            description = $2,
            amount = $3,
            type = $4,
            category = $5,
            event_date = $6
        WHERE id = $7;
    `
	_, err := p.DB.ExecContext(
		ctx, query,
		t.Name, t.Description, t.Amount, t.Type, t.Category, t.EventDate,
		t.ID,
	)
	return err
}

func (p *Postgres) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM transactions WHERE id = $1;`
	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}

func (p *Postgres) GetByPeriod(ctx context.Context, from, to time.Time) ([]model.Transaction, error) {
	query := `
	SELECT * FROM transactions 
	WHERE event_date 
	BETWEEN $1 AND $2 
	ORDER BY event_date
	`
	var t []model.Transaction
	err := p.DB.SelectContext(ctx, &t, query, from, to)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (p *Postgres) GetAllSorted(ctx context.Context, sortField, order string) ([]model.Transaction, error) {
	if !sortFields[sortField] {
		sortField = "id"
	}
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	var t []model.Transaction
	query := fmt.Sprintf(`SELECT * FROM transactions ORDER BY %s %s`, sortField, order)
	err := p.DB.SelectContext(ctx, &t, query)
	if err != nil {
		return nil, err
	}
	return t, nil

}

func (p *Postgres) GetSum(ctx context.Context, from, to time.Time) (int64, error) {
	query := `
        SELECT COALESCE(SUM(amount), 0)
        FROM transactions
        WHERE event_date BETWEEN $1 AND $2;
    `
	var sum int64
	err := p.DB.GetContext(ctx, &sum, query, from, to)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (p *Postgres) GetAvg(ctx context.Context, from, to time.Time) (float64, error) {
	query := `
        SELECT COALESCE(AVG(amount), 0)
        FROM transactions
        WHERE event_date BETWEEN $1 AND $2;
    `
	var avg float64
	err := p.DB.GetContext(ctx, &avg, query, from, to)
	if err != nil {
		return 0, err
	}
	return avg, nil
}

func (p *Postgres) GetCount(ctx context.Context, from, to time.Time) (int64, error) {
	query := `
        SELECT COUNT(*)
        FROM transactions
        WHERE event_date BETWEEN $1 AND $2;
    `
	var count int64
	err := p.DB.GetContext(ctx, &count, query, from, to)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *Postgres) GetMedian(ctx context.Context, from, to time.Time) (float64, error) {
	query := `
        SELECT COALESCE(
            percentile_cont(0.5) WITHIN GROUP (ORDER BY amount),
            0
        )
        FROM transactions
        WHERE event_date BETWEEN $1 AND $2;
    `
	var median float64
	err := p.DB.GetContext(ctx, &median, query, from, to)
	if err != nil {
		return 0, err
	}
	return median, nil
}

func (p *Postgres) GetPercentile90(ctx context.Context, from, to time.Time) (float64, error) {
	query := `
        SELECT COALESCE(
            percentile_cont(0.9) WITHIN GROUP (ORDER BY amount),
            0
        )
        FROM transactions
        WHERE event_date BETWEEN $1 AND $2;
    `
	var result float64
	err := p.DB.GetContext(ctx, &result, query, from, to)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (p *Postgres) GroupByDay(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	query := `
		SELECT to_char(event_date,'YYYY-MM-DD') AS day, SUM(amount) AS total 
		FROM transactions 
		WHERE event_date 
		BETWEEN $1 AND $2 
		GROUP BY day 
		ORDER BY day;
	`
	rows, err := p.DB.QueryxContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var day string
		var total int64
		err = rows.Scan(&day, &total)
		if err != nil {
			return nil, err
		}
		result[day] = total

	}
	return result, nil
}

func (p *Postgres) GroupByWeek(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	query := `
		SELECT to_char(date_trunc('week', event_date),'IYYY-IW') AS week, SUM(amount) AS total 
		FROM transactions 
		WHERE event_date 
		BETWEEN $1 AND $2 
		GROUP BY week 
		ORDER BY week;
	`
	rows, err := p.DB.QueryxContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var week string
		var total int64
		err = rows.Scan(&week, &total)
		if err != nil {
			return nil, err
		}
		result[week] = total
	}
	return result, nil

}

func (p *Postgres) GroupByMonth(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	query := `
			SELECT to_char(date_trunc('month', event_date),'YYYY-MM') AS month, SUM(amount) AS total 
			FROM transactions 
			WHERE event_date 
			BETWEEN $1 AND $2 
			GROUP BY month 
			ORDER BY month;
			`
	rows, err := p.DB.QueryxContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var month string
		var total int64
		if err := rows.Scan(&month, &total); err != nil {
			return nil, err
		}
		result[month] = total
	}
	return result, nil

}

func (p *Postgres) GroupByCategory(ctx context.Context, from, to time.Time) (map[string]int64, error) {
	query := `
		SELECT category, SUM(amount) AS total 
		FROM transactions 
		WHERE event_date 
		BETWEEN $1 AND $2 
		GROUP BY category 
		ORDER BY category;
	`
	rows, err := p.DB.QueryxContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var category string
		var total int64
		err = rows.Scan(&category, &total)
		if err != nil {
			return nil, err
		}
		result[category] = total
	}
	return result, nil

}
