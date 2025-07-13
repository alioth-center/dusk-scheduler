package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/app/repository"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
	"io"
	"net/http"
	"net/url"
	"time"
)

type outComeService struct {
	taskDao    repository.TaskDao
	painterDao repository.PainterDao
	outcomeDao repository.OutcomeDao
	storageDao repository.StorageDao
	sysLogger  logger.Logger
	client     *http.Client
}

func (srv *outComeService) CreateOutcome(ctx context.Context, painterName string, taskID uint64, reference string, createdAt, completedAt time.Time) (err error) {
	task, existTask, getTaskErr := srv.taskDao.GetTaskByID(ctx, taskID)
	if getTaskErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get task by id: %d", taskID), getTaskErr)

		return getTaskErr
	}
	if !existTask {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("task does not exist: %d", taskID), nil)

		return errors.TaskReferenceNotFoundError()
	}

	painter, existPainter, getPainterErr := srv.painterDao.GetPainterByName(ctx, painterName)
	if getPainterErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get painter by name: %s", painterName), getPainterErr)

		return getPainterErr
	}
	if !existPainter {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("painter does not exist: %s", painterName), nil)

		return errors.GetOutcomeContentPainterNotFoundError()
	}

	outcomeEntity := domain.Outcome{
		Instance:    painter.ID,
		TaskID:      task.ID,
		Reference:   reference,
		StartedAt:   createdAt,
		CompletedAt: completedAt,
	}
	if _, createErr := srv.outcomeDao.CreateOutcome(ctx, &outcomeEntity); createErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to create outcome: %v", createErr), createErr)

		return createErr
	}

	return nil
}

func (srv *outComeService) GetOutcomeByTaskID(ctx context.Context, taskID uint64) (outcome *domain.Outcome, exist bool, err error) {
	outcome, exist, err = srv.outcomeDao.GetOutcomeByTaskID(ctx, taskID)
	if err != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get outcome by task id: %d", taskID), err)

		return nil, false, err
	}

	return outcome, exist, nil
}

func (srv *outComeService) GetOutcomeByReference(ctx context.Context, reference string) (outcome *domain.Outcome, exist bool, err error) {
	outcome, exist, err = srv.outcomeDao.GetOutcomeByReference(ctx, reference)
	if err != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get outcome by reference: %s", reference), err)

		return nil, false, err
	}

	return outcome, exist, nil
}

func (srv *outComeService) GetOutcomeContent(ctx context.Context, reference string) (content *bytes.Buffer, err error) {
	contentURL, queryErr := srv.GetOutcomeURL(ctx, reference)
	if queryErr != nil {
		return nil, queryErr
	}

	request, buildRequestErr := http.NewRequestWithContext(ctx, http.MethodGet, contentURL.String(), nil)
	if buildRequestErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to build request for content: %s", contentURL.String()), buildRequestErr)

		return nil, buildRequestErr
	}

	response, executeErr := srv.client.Do(request)
	if executeErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to execute request for content: %s", contentURL.String()), executeErr)

		return nil, executeErr
	}

	defer errors.Ignore(response.Body.Close())
	if response.StatusCode != http.StatusOK {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("failed to get outcome content from content: %s(%d)", contentURL.String(), response.StatusCode), nil)

		return nil, errors.GetOutcomeContentDownloadFailedError()
	}

	content = &bytes.Buffer{}
	if _, copyErr := io.Copy(content, response.Body); copyErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to copy content: %s", contentURL.String()), copyErr)

		return nil, copyErr
	}

	return content, nil
}

func (srv *outComeService) GetOutcomeURL(ctx context.Context, reference string) (content *url.URL, err error) {
	outcome, existOutcome, getOutcomeErr := srv.outcomeDao.GetOutcomeByReference(ctx, reference)
	if getOutcomeErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get outcome by reference: %s", reference), err)

		return nil, getOutcomeErr
	}
	if !existOutcome {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("outcome does not exist: %s", reference), nil)

		return nil, errors.GetOutcomeContentOutcomeNotFoundError()
	}

	painter, existPainter, getPainterErr := srv.painterDao.GetPainterByID(ctx, outcome.Instance)
	if getPainterErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get painter by id: %d", outcome.Instance), err)

		return nil, getPainterErr
	}
	if !existPainter {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("painter does not exist: %s", outcome.Instance), nil)

		return nil, errors.GetOutcomeContentPainterNotFoundError()
	}

	storage, existStorage, getStorageErr := srv.storageDao.GetStorageByID(ctx, painter.PolicyID)
	if getStorageErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to get storage by id: %d", painter.PolicyID), err)

		return nil, getStorageErr
	}
	if !existStorage {
		srv.sysLogger.WarnCtx(ctx, fmt.Sprintf("storage does not exist: %s", painter.PolicyID), nil)

		return nil, errors.GetOutcomeContentStorageNotFoundError()
	}

	outcomeURL := fmt.Sprintf(storage.Formatter, outcome.Reference)
	parsedURL, parseErr := url.Parse(outcomeURL)
	if parseErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to parse url: %s", outcomeURL), err)

		return nil, parseErr
	}

	return parsedURL, nil
}
