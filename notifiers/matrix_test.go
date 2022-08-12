package notifiers

import (
	"testing"
	"time"

	"github.com/statping-ng/statping-ng/database"
	"github.com/statping-ng/statping-ng/types/core"
	"github.com/statping-ng/statping-ng/types/failures"
	"github.com/statping-ng/statping-ng/types/notifications"
	"github.com/statping-ng/statping-ng/types/null"
	"github.com/statping-ng/statping-ng/types/services"
	"github.com/statping-ng/statping-ng/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	matrixHost   string
	matrixRoomId string
	matrixToken  string
)

func TestMatrixNotifier(t *testing.T) {
	err := utils.InitLogs()
	require.Nil(t, err)

	t.Parallel()

	matrixHost = utils.Params.GetString("MATRIX_HOST")
	matrixRoomId = utils.Params.GetString("MATRIX_ROOM_ID")
	matrixToken = utils.Params.GetString("MATRIX_TOKEN")
	if matrixHost == "" || matrixRoomId == "" || matrixToken == "" {
		t.Log("Matrix notifier testing skipped, missing MATRIX_HOST, MATRIX_ROOM_ID and MATRIX_TOKEN environment variable")
		t.SkipNow()
	}

	Matrix.Host = null.NewNullString(matrixHost)
	Matrix.Var1 = null.NewNullString(matrixRoomId)
	Matrix.ApiSecret = null.NewNullString(matrixToken)

	db, err := database.OpenTester()
	require.Nil(t, err)
	db.AutoMigrate(&notifications.Notification{})
	notifications.SetDB(db)
	core.Example()

	t.Run("Load Matrix", func(t *testing.T) {
		Matrix.ApiSecret = null.NewNullString(matrixToken)
		Matrix.Var1 = null.NewNullString(matrixRoomId)
		Matrix.Delay = time.Duration(1 * time.Second)
		Matrix.Enabled = null.NewNullBool(true)

		Add(Matrix)

		assert.Equal(t, "jojo", Matrix.Author)
		assert.Equal(t, matrixToken, Matrix.ApiSecret.String)
		assert.Equal(t, matrixRoomId, Matrix.Var1.String)
		assert.Equal(t, matrixHost, Matrix.Host.String)
	})

	t.Run("Matrix Within Limits", func(t *testing.T) {
		assert.True(t, Matrix.CanSend())
	})

	t.Run("Matrix OnSave", func(t *testing.T) {
		_, err := Matrix.OnSave()
		assert.Nil(t, err)
	})

	t.Run("Matrix OnFailure", func(t *testing.T) {
		_, err := Matrix.OnFailure(services.Example(false), failures.Example())
		assert.Nil(t, err)
	})

	t.Run("Matrix OnSuccess", func(t *testing.T) {
		_, err := Matrix.OnSuccess(services.Example(true))
		assert.Nil(t, err)
	})

	t.Run("Matrix Test", func(t *testing.T) {
		_, err := Matrix.OnTest()
		assert.Nil(t, err)
	})
}
