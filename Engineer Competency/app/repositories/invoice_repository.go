package repositories

import (
	"engineer-comp/app/global/pgsql"
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/models"
	"strconv"

	"strings"

	"fmt"
	"time"
)

type InvoiceRepositoryInterface interface {
	CreateInvoice(request *models.InvoiceRequest, resultChan chan *models.InvoiceResponseChan)
	GetInvoiceByID(id int, resultChan chan *models.InvoiceResponseChan)
	GetAllInvoice(resultChan chan *models.InvoiceListResponseChan)
	UpdateInvoice(request *models.InvoiceRequest, resultChan chan *models.InvoiceResponseChan)
	DeleteInvoice(id int, resultChan chan *models.InvoiceResponseChan)
}

type invoiceRepository struct {
	sqldb pgsql.SqlInterface
}

func InitInvoiceRepository(sqldb pgsql.SqlInterface) InvoiceRepositoryInterface {
	return &invoiceRepository{
		sqldb: sqldb,
	}
}

// done
func (r *invoiceRepository) CreateInvoice(request *models.InvoiceRequest, resultChan chan *models.InvoiceResponseChan) {

	response := &models.InvoiceResponseChan{}

	query := `INSERT INTO public.invoice
	(
		user_id,
		merchant_id,
		invoice_name,
		amount,
		status,
		created_at
	)
	VALUES
	(
		$1, $2, $3, $4, 'Pending', NOW()
	)
	RETURNING id`

	var invoiceID int

	err := r.sqldb.DB().QueryRow(
		query,
		request.BuyerId,
		request.MerchantID,
		request.InvoiceName,
		request.Amount,
	).Scan(&invoiceID)

	if err != nil {
		fmt.Println("error create invoice")
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	request.ID = invoiceID

	response.Invoice = &models.Invoice{
		ID:          request.ID,
		BuyerId:     request.BuyerId,
		MerchantID:  request.MerchantID,
		InvoiceName: request.InvoiceName,
		Amount:      request.Amount,
		Status:      request.Status,
	}
	response.Error = nil

	resultChan <- response
}

// done
func (r *invoiceRepository) GetInvoiceByID(id int, resultChan chan *models.InvoiceResponseChan) {

	response := &models.InvoiceResponseChan{}

	var total int

	sqlCount := `SELECT COUNT(id) FROM public.invoice WHERE id = $1 AND deleted_at IS NULL`

	err := r.sqldb.DB().QueryRow(sqlCount, id).Scan(&total)

	if err != nil {
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	if total == 0 {
		response.Error = helper.NewError("invoice not found")
		resultChan <- response
		return
	}

	query := `SELECT
		id,
		user_id,
		merchant_id,
		COALESCE(admin_id::bigint, '0') AS admin_id,
		invoice_name,
		amount,
		status,
		COALESCE(payment_type::text, '-') AS payment_type,
		created_at
	FROM public.invoice
	WHERE id = $1
	AND deleted_at IS NULL`

	rows, err := r.sqldb.DB().Query(query, id)

	if err != nil {
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	defer rows.Close()

	var invoice models.Invoice

	for rows.Next() {

		err = rows.Scan(
			&invoice.ID,
			&invoice.BuyerId,
			&invoice.MerchantID,
			&invoice.AdminID,
			&invoice.InvoiceName,
			&invoice.Amount,
			&invoice.Status,
			&invoice.PaymentType,
			&invoice.CreatedAt,
		)

		if err != nil {
			fmt.Println(err.Error())

			response.Error = err
			resultChan <- response
			return
		}
	}

	response.Invoice = &invoice
	response.Error = nil

	resultChan <- response
}

// done
func (r *invoiceRepository) GetAllInvoice(resultChan chan *models.InvoiceListResponseChan) {

	response := &models.InvoiceListResponseChan{}

	query := `SELECT
		id,
		user_id,
		merchant_id,
		COALESCE(admin_id::bigint, '0') AS admin_id,
		invoice_name,
		amount,
		status,
		COALESCE(payment_type::text, '-') AS payment_type,
		created_at
	FROM public.invoice
	WHERE deleted_at IS NULL
	ORDER BY id DESC`

	rows, err := r.sqldb.DB().Query(query)

	if err != nil {
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	defer rows.Close()

	var invoices []models.Invoice

	for rows.Next() {

		var invoice models.Invoice

		err = rows.Scan(
			&invoice.ID,
			&invoice.BuyerId,
			&invoice.MerchantID,
			&invoice.AdminID,
			&invoice.InvoiceName,
			&invoice.Amount,
			&invoice.Status,
			&invoice.PaymentType,
			&invoice.CreatedAt,
		)

		if err != nil {
			fmt.Println(err.Error())

			response.Error = err
			resultChan <- response
			return
		}

		invoices = append(invoices, invoice)
	}

	response.Invoice = invoices
	response.Error = nil

	resultChan <- response
}

func (r *invoiceRepository) UpdateInvoice(request *models.InvoiceRequest, resultChan chan *models.InvoiceResponseChan) {

	response := &models.InvoiceResponseChan{}
	var query string

	walletRepo := &walletRepository{
		sqldb: r.sqldb,
	}
	//--status (pending, paid, approved, rejected, refund)
	if request.RoleID == "1" {
		//payment
		if request.Status == "paid" {
			if strings.ToLower(request.PaymentType) == "wallet" {
				// cek balance
				balance, err := walletRepo.GetBalanceByUserID(request.BuyerId)
				if err != nil {
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}

				//get bill
				bill, err := r.GetInvoiceBillByID(request.ID)
				if err != nil {
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}

				if balance < bill {
					err = helper.NewError("balance wallet tidak cukup dengan tagihan")
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}
				//paid and update wallet

				if request.Status == "paid" {
					query = `BEGIN;
								INSERT INTO public.wallet
									( user_id, invoice_id, current_saldo, record_saldo, "action",status, created_at)
									select i.user_id, 
											i.id, 
											CASE
												WHEN i.status = 'pending' THEN  w.current_saldo - i.amount
												ELSE w.current_saldo 
											END AS current_saldo,
											i.amount,
											'wallet payment',
											1,
											now()
									from invoice i
									join (select * from wallet  where user_id = ` + strconv.Itoa(int(request.BuyerId)) + ` order by id desc limit 1) w on i.user_id = w.user_id 
									where i.id = ` + strconv.Itoa(request.ID) + ` ; 
								UPDATE public.invoice
									SET
										status = 'paid',
										payment_type = 'wallet',
										updated_at = now()
									WHERE id = ` + strconv.Itoa(request.ID) + `
									AND deleted_at IS NULL;
							COMMIT;`
					_, err := r.sqldb.DB().Exec(query)
					if err != nil {
						fmt.Println(err.Error())
						response.Error = err
						resultChan <- response
						return
					}
				} else {
					err := helper.NewError("invalid invoice status")
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}
			} else if strings.ToLower(request.PaymentType) == "va" || strings.ToLower(request.PaymentType) == "ewallet" {
				balance := request.VAWallet
				//get bill
				bill, err := r.GetInvoiceBillByID(request.ID)
				if err != nil {
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}

				if balance != bill {
					err = helper.NewError("payment ewallet / va tidak sesuai dengan tagihan")
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}

				query = `
						UPDATE public.invoice
							SET
								status = 'paid',
								payment_type = 'Ewallet / VA',
								updated_at = now()
							WHERE id = $1
							AND deleted_at IS NULL;
					`
				_, err = r.sqldb.DB().Exec(
					query,
					request.ID,
				)
				if err != nil {
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}

			} else {
				err := helper.NewError("invalid payment type")
				fmt.Println(err.Error())
				response.Error = err
				resultChan <- response
				return
			}
		} else if request.Status == "refund" {
			getInvoiceChan := make(chan *models.InvoiceResponseChan)
			go r.GetInvoiceByID(request.ID, getInvoiceChan)
			invoiceResult := <-getInvoiceChan

			if invoiceResult.Error != nil {
				//fail get invoice by id
				fmt.Println(invoiceResult.Error.Error())
				response.Error = invoiceResult.Error
				resultChan <- response
				return
			}

			if invoiceResult.Invoice.Status == "approved" {
				query = `BEGIN;
						UPDATE public.invoice
							SET
								status = 'refund',
								updated_at = now()
							WHERE id = ` + strconv.Itoa(request.ID) + ` and user_id = ` + strconv.Itoa(int(request.BuyerId)) + `
							AND deleted_at IS NULL;
						COMMIT;`
				_, err := r.sqldb.DB().Exec(query)
				if err != nil {
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}
			}
		} else {
			err := helper.NewError("invalid invoice status")
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}

	} else if request.RoleID == "3" {
		getInvoiceChan := make(chan *models.InvoiceResponseChan)
		go r.GetInvoiceByID(request.ID, getInvoiceChan)
		invoiceResult := <-getInvoiceChan

		if invoiceResult.Error != nil {
			//fail get invoice by id
			fmt.Println(invoiceResult.Error.Error())
			response.Error = invoiceResult.Error
			resultChan <- response
			return
		}
		fmt.Println(request.Status)
		//approve
		if request.Status == "approved" {
			if invoiceResult.Invoice.Status == "paid" {
				query = `BEGIN;
						INSERT INTO public.wallet
							( user_id, invoice_id, current_saldo, record_saldo, "action", status, created_at)
							select i.merchant_id  , 
									i.id, 
									CASE
										WHEN i.status = 'paid' THEN i.amount + w.current_saldo 
										ELSE w.current_saldo 
									END AS current_saldo,
									i.amount,
									'income',
									1,
									now()
							from invoice i
							join (select * from wallet where user_id = ` + strconv.Itoa(int(invoiceResult.Invoice.MerchantID)) + ` order by created_at desc limit 1) w on i.merchant_id  = w.user_id 
							where i.id = ` + strconv.Itoa(request.ID) + ` ; 
						UPDATE public.invoice
							SET
								admin_id = ` + strconv.Itoa(int(request.AdminID)) + `,
								status = 'approved',
								updated_at = now()
							WHERE id = ` + strconv.Itoa(request.ID) + `
							AND deleted_at IS NULL;
					COMMIT;`
				_, err := r.sqldb.DB().Exec(query)
				if err != nil {
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}
			} else {
				err := helper.NewError("only paid invoice can be approved")
				fmt.Println(err.Error())
				response.Error = err
				resultChan <- response
				return
			}
		} else if request.Status == "rejected" {
			if invoiceResult.Invoice.Status == "paid" {
				query = `BEGIN;
						INSERT INTO public.wallet
							( user_id, invoice_id, current_saldo, record_saldo, "action", created_at)
							select i.user_id  , 
									i.id, 
									CASE
										WHEN i.status = 'paid' THEN i.amount + w.current_saldo 
										ELSE w.current_saldo 
									END AS current_saldo,
									i.amount,
									'refund transaction',
									now()
							from invoice i
							join (select * from wallet where user_id = ` + strconv.Itoa(int(invoiceResult.Invoice.BuyerId)) + ` order by created_at desc limit 1 ) w on i.user_id  = w.user_id 
							where i.id = ` + strconv.Itoa(request.ID) + ` ; 

						UPDATE public.invoice
							SET
								admin_id = ` + strconv.Itoa(int(request.AdminID)) + `,
								status = 'rejected',
								updated_at = now()
							WHERE id = ` + strconv.Itoa(request.ID) + `
							AND deleted_at IS NULL;
					COMMIT;`
				_, err := r.sqldb.DB().Exec(query)
				if err != nil {
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}
			} else {
				err := helper.NewError("only paid invoice can be rejected")
				fmt.Println(err.Error())
				response.Error = err
				resultChan <- response
				return
			}

		} else if request.Status == "refund" {
			fmt.Println("masuk refund")
			if invoiceResult.Invoice.Status == "refund" {
				query = `BEGIN;
						INSERT INTO public.wallet
							( user_id, invoice_id, current_saldo, record_saldo, "action", created_at)
							select i.user_id  , 
									i.id, 
									CASE
										WHEN i.status = 'refund' THEN i.amount + w.current_saldo 
										ELSE w.current_saldo 
									END AS current_saldo,
									i.amount,
									'refund transaction',
									now()
							from invoice i
							join (select * from wallet where user_id = ` + strconv.Itoa(int(invoiceResult.Invoice.BuyerId)) + ` order by id desc limit 1 ) w on i.user_id  = w.user_id 
							where i.id = ` + strconv.Itoa(request.ID) + ` ; 

						INSERT INTO public.wallet
						( user_id, invoice_id, current_saldo, record_saldo, "action", created_at)
						select i.merchant_id  , 
								i.id, 
								CASE
									WHEN i.status = 'refund' THEN  w.current_saldo - i.amount 
									ELSE w.current_saldo 
								END AS current_saldo,
								i.amount,
								'refund transaction',
								now()
						from invoice i
						join (select * from wallet where user_id = ` + strconv.Itoa(int(invoiceResult.Invoice.MerchantID)) + ` order by id desc limit 1) w on i.merchant_id  = w.user_id 
						where i.id = ` + strconv.Itoa(request.ID) + ` ; 


						UPDATE public.invoice
							SET
								admin_id = ` + strconv.Itoa(int(request.AdminID)) + `,
								status = 'refunded',
								updated_at = now()
							WHERE id = ` + strconv.Itoa(request.ID) + `
							AND deleted_at IS NULL;
					COMMIT;`
				_, err := r.sqldb.DB().Exec(query)
				if err != nil {
					fmt.Println(err.Error())
					response.Error = err
					resultChan <- response
					return
				}
			} else {
				err := helper.NewError("only refund invoice can be refunded")
				fmt.Println(err.Error())
				response.Error = err
				resultChan <- response
				return
			}

		} else {
			err := helper.NewError("invalid invoice status")
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}
	}

	// status

	response.Invoice = &models.Invoice{
		ID:          request.ID,
		BuyerId:     request.BuyerId,
		MerchantID:  request.MerchantID,
		AdminID:     request.AdminID,
		InvoiceName: request.InvoiceName,
		Amount:      request.Amount,
		Status:      request.Status,
	}
	response.Error = nil

	resultChan <- response
}

func (r *invoiceRepository) DeleteInvoice(id int, resultChan chan *models.InvoiceResponseChan) {

	response := &models.InvoiceResponseChan{}

	query := `UPDATE public.invoice
	SET deleted_at = $1
	WHERE id = $2`

	_, err := r.sqldb.DB().Exec(
		query,
		time.Now(),
		id,
	)

	if err != nil {
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	response.Error = nil

	resultChan <- response
}

func (r *invoiceRepository) PayInvoice(request *models.Invoice, resultChan chan *models.InvoiceResponseChan) {

	response := &models.InvoiceResponseChan{}

	query := `UPDATE public.invoice
	SET
		admin_id = $1,
		status = $2,
		updated_at = $3
	WHERE id = $4
	AND deleted_at IS NULL`

	_, err := r.sqldb.DB().Exec(
		query,
		request.AdminID,
		request.Status,
		time.Now(),
		request.ID,
	)

	if err != nil {
		fmt.Println(err.Error())

		response.Error = err
		resultChan <- response
		return
	}

	// status

	response.Invoice = request
	response.Error = nil

	resultChan <- response
}

func (r *invoiceRepository) GetInvoiceBillByID(id int) (float64, error) {

	var total int

	sqlCount := `SELECT COUNT(id) FROM public.invoice WHERE id = $1 AND deleted_at IS NULL`

	err := r.sqldb.DB().QueryRow(sqlCount, id).Scan(&total)

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	if total == 0 {
		return 0, helper.NewError("invoice not found")
	}

	query := `SELECT
				amount
			FROM public.invoice
			WHERE id = $1
			AND deleted_at IS NULL`

	rows, err := r.sqldb.DB().Query(query, id)

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	defer rows.Close()

	var invoiceBill float64

	for rows.Next() {

		err = rows.Scan(
			&invoiceBill,
		)

		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}
	}
	return invoiceBill, nil
}
