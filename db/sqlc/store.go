package db

import (
    "context"
    "database/sql"
    "fmt"
)

// Store provides all functions to execute SQL queries and transactions
type Store interface {
    Querier
    TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
    db *sql.DB
    *Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
    return &SQLStore{
        db:      db,
        Queries: New(db),
    }
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
    tx, err := store.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    q := New(tx)
    err = fn(q)
    if err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
        }
        return err
    }

    return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
    FromAccountID int64 `json:"from_account_id"`
    ToAccountID   int64 `json:"to_account_id"`
    Amount        int64 `json:"amount"`
}

// TransferTxResult  is the result of the transfer transaction
type TransferTxResult struct {
    Transfer    Transfer `json:"transfer"`
    FromAccount Account  `json:"from_account"`
    ToAccount   Account  `json:"to_account"`
    FromEntry   Entry    `json:"from_entry"`
    ToEntry     Entry    `json:"to_entry"`
}

// untuk mengatasi deadlock, maka context.Background() di go func store_test.go  diganti
// sehingga membutuhkan key {Empty Struct}
// var transactionKey = struct{}{}

// Perform a money  transfer from one account to the other
// it creates a transfer record, add  account entries, and update accounts balance within a single database transacction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
    var result TransferTxResult

    err := store.execTx(ctx, func(q *Queries) error {

        var err error

        // transactionName := ctx.Value(transactionKey)

        // fmt.Println(transactionName, "create transfer")
        result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
            FromAccountID: arg.FromAccountID,
            ToAccountID:   arg.ToAccountID,
            Amount:        arg.Amount,
        })
        if err != nil {
            return err
        }

        // fmt.Println(transactionName, "create entry 1")
        result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.FromAccountID,
            Amount:    -arg.Amount,
        })
        if err != nil {
            return err
        }

        // fmt.Println(transactionName, "create entry 2")
        result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.ToAccountID,
            Amount:    arg.Amount,
        })
        if err != nil {
            return err
        }

        /* OLD WAY
        // update account's 1 balance
        // fmt.Println(transactionName, "get account 1")
        account1, err  := q.GetAccountForUpdate(ctx, arg.FromAccountID)
        if err != nil {
            return err
        }

        // fmt.Println(transactionName, "update account 1")
        result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
            ID: arg.FromAccountID,
            Balance: account1.Balance - arg.Amount,
        })

        // update account's 2 balance
        // fmt.Println(transactionName, "get account 2")
        account2, err  := q.GetAccountForUpdate(ctx, arg.ToAccountID)
        if err != nil {
            return err
        }

        // fmt.Println(transactionName, "update account 2")
        result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
            ID: arg.ToAccountID,
            Balance: account2.Balance + arg.Amount,
        })

        */
        /*
            1st NEW WAY
            result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
                ID:     arg.FromAccountID,
                Amount: -arg.Amount,
            })
            if err != nil {
                return err
            }

            result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
                ID:     arg.ToAccountID,
                Amount: arg.Amount,
            })
            if err != nil {
                return err
            }
        */

        //2nd Way, to avoid deadlock
        if arg.FromAccountID < arg.ToAccountID {
            result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
        } else {
            result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
        }

        return nil
    })

    return result, err
}

func addMoney(
    ctx context.Context,
    q *Queries,
    accountID1 int64,
    amount1 int64,
    accountID2 int64,
    amount2 int64,
) (account1 Account, account2 Account, err error) {
    account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        ID:     accountID1,
        Amount: amount1,
    })
    if err != nil {
        return
    }

    account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        ID:     accountID2,
        Amount: amount2,
    })
    return
}