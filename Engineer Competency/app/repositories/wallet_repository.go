package repositories

import (
	"engineer-comp/app/global/pgsql"
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/models"
	"strconv"

	"fmt"
)

type WalletRepositoryInterface interface {
	CreateWallet(request *models.Wallet, resultChan chan *models.WalletResponseChan)
	TopUpWallet(request *models.Wallet, resultChan chan *models.WalletResponseChan)
	GetWalletByID(id int, resultChan chan *models.WalletResponseChan)
	GetWalletByUserID(userID int, resultChan chan *models.WalletListResponseChan)
	GetBalanceByUserID(userID int64) (float64, error)
	GetBalanceWalletByUserID(userID int, resultChan chan *models.WalletResponseChan)
}

type walletRepository struct {
	sqldb pgsql.SqlInterface
}

func InitWalletRepository(
	sqldb pgsql.SqlInterface,
) WalletRepositoryInterface {

	return &walletRepository{
		sqldb: sqldb,
	}
}

func (r *walletRepository) CreateWallet(request *models.Wallet, resultChan chan *models.WalletResponseChan) {

	response := &models.WalletResponseChan{}

	sql := `
		INSERT INTO public.wallet
		(
			user_id,
			invoice_id,
			current_saldo,
			record_saldo,
			action,
			status,
			created_at
		)
		VALUES
		(
			$1,
			0,
			0,
			0,
			'Saldo Awal',
			1,
			NOW()
		)
		RETURNING id
	`

	err := r.sqldb.DB().QueryRow(
		sql,
		request.UserID,
	).Scan(&request.ID)

	if err != nil {
		fmt.Println("create wallet error")
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	response.Wallet = request
	response.Error = nil

	resultChan <- response
	return
}

func (r *walletRepository) TopUpWallet(request *models.Wallet, resultChan chan *models.WalletResponseChan) {

	response := &models.WalletResponseChan{}
	fmt.Println(request.RecordSaldo)
	sql := `
		insert into wallet 
			( user_id, invoice_id, current_saldo, record_saldo, "action", status, created_at)
		select $1, 0, x.current_saldo + $2, $3, 'Top Up', 1, now()  from 
			(select * from wallet where user_id = $4 order by id desc limit 1) x
		RETURNING id
	`

	err := r.sqldb.DB().QueryRow(
		sql,
		request.UserID,
		request.RecordSaldo,
		request.RecordSaldo,
		request.UserID,
	).Scan(&request.ID)

	if err != nil {
		fmt.Println("Top up wallet error")
		fmt.Println(err.Error())
		response.Error = err
		resultChan <- response
		return
	}

	if err != nil {
		fmt.Println("create wallet error")
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	response.Wallet = request
	response.Error = nil

	resultChan <- response
	return
}

func (r *walletRepository) GetWalletByID(id int, resultChan chan *models.WalletResponseChan) {

	response := &models.WalletResponseChan{}

	var total int

	sqlCount := `
		SELECT COUNT(id)
		FROM public.wallet
		WHERE id = $1
		AND deleted_at IS NULL
	`

	err := r.sqldb.DB().QueryRow(sqlCount, id).Scan(&total)

	if err != nil {
		fmt.Println("get wallet by id count error")
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	if total == 0 {
		err = helper.NewError("wallet not found")

		response.Error = err
		resultChan <- response
		return
	}

	sql := `
		SELECT
			id,
			user_id,
			invoice_id,
			current_saldo,
			record_saldo,
			action,
			status,
			created_at,
			updated_at
		FROM public.wallet
		WHERE id = $1
		AND deleted_at IS NULL
	`

	rows, err := r.sqldb.DB().Query(sql, id)

	if err != nil {
		fmt.Println("get wallet by id query error")
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	defer rows.Close()

	var wallet models.Wallet

	for rows.Next() {

		err = rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.InvoiceID,
			&wallet.CurrentSaldo,
			&wallet.RecordSaldo,
			&wallet.Action,
			&wallet.Status,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
		)

		if err != nil {
			fmt.Println("scan wallet error")
			fmt.Println(err.Error())

			response.Error = err
			resultChan <- response
			return
		}
	}

	response.Wallet = &wallet
	response.Error = nil

	resultChan <- response
	return
}

func (r *walletRepository) GetWalletByUserID(userID int, resultChan chan *models.WalletListResponseChan) {

	response := &models.WalletListResponseChan{}

	sql := `
		SELECT
			id,
			user_id,
			invoice_id,
			current_saldo,
			record_saldo,
			action,
			status,
			created_at,
			updated_at
		FROM public.wallet
		WHERE user_id = $1
		AND deleted_at IS NULL
		ORDER BY id DESC
	`

	rows, err := r.sqldb.DB().Query(sql, userID)

	if err != nil {
		fmt.Println("get wallet by user id error")
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	defer rows.Close()

	var wallets []models.Wallet

	for rows.Next() {

		var wallet models.Wallet

		err = rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.InvoiceID,
			&wallet.CurrentSaldo,
			&wallet.RecordSaldo,
			&wallet.Action,
			&wallet.Status,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
		)

		if err != nil {
			fmt.Println("scan wallet list error")
			fmt.Println(err.Error())

			response.Error = err
			resultChan <- response
			return
		}

		wallets = append(wallets, wallet)
	}

	response.Wallet = wallets
	response.Error = nil

	resultChan <- response
	return
}

func (r *walletRepository) GetBalanceWalletByUserID(userID int, resultChan chan *models.WalletResponseChan) {

	response := &models.WalletResponseChan{}

	sql := `select current_saldo from wallet where user_id = $1 order by id desc limit 1`

	rows, err := r.sqldb.DB().Query(sql, userID)

	if err != nil {
		fmt.Println("get wallet by user id error")
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	defer rows.Close()

	var wallet models.Wallet

	for rows.Next() {

		err = rows.Scan(
			&wallet.CurrentSaldo,
		)

		if err != nil {
			fmt.Println("scan wallet list error")
			fmt.Println(err.Error())

			response.Error = err
			resultChan <- response
			return
		}

	}

	response.Wallet = &wallet
	response.Error = nil

	resultChan <- response
	return
}

func (r *walletRepository) GetBalanceByUserID(userID int64) (float64, error) {
	sql := `select coalesce(current_saldo, 0) from wallet where user_id = ` + strconv.FormatInt(userID, 10) + ` order by id desc limit 1`
	rows, err := r.sqldb.DB().Query(sql)
	if err != nil {
		fmt.Println("get wallet by user id error")
		fmt.Println(err.Error())
		return 0, err
	}

	defer rows.Close()

	var wallet models.Wallet

	for rows.Next() {

		err = rows.Scan(
			&wallet.CurrentSaldo,
		)

		if err != nil {
			fmt.Println("scan wallet list error")
			fmt.Println(err.Error())
			return 0, err
		}

	}

	return wallet.CurrentSaldo, nil
}
