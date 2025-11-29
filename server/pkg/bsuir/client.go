package bsuir

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"schedluer/internal/config"
	"schedluer/internal/models"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg *config.BSUIRAPIConfig) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

func (c *Client) GetGroupSchedule(groupNumber string) (*models.ScheduleResponse, error) {
	url := fmt.Sprintf("%s/schedule?studentGroup=%s", c.baseURL, groupNumber)

	var response models.ScheduleResponse
	if err := c.makeRequest(url, &response); err != nil {
		return nil, fmt.Errorf("failed to get group schedule: %w", err)
	}

	return &response, nil
}

func (c *Client) GetEmployeeSchedule(urlID string) (*models.ScheduleResponse, error) {
	url := fmt.Sprintf("%s/employees/schedule/%s", c.baseURL, urlID)

	var response models.ScheduleResponse
	if err := c.makeRequest(url, &response); err != nil {
		return nil, fmt.Errorf("failed to get employee schedule: %w", err)
	}

	return &response, nil
}

func (c *Client) GetAllGroups() ([]models.StudentGroupListItem, error) {
	url := fmt.Sprintf("%s/student-groups", c.baseURL)

	var groups []models.StudentGroupListItem
	if err := c.makeRequest(url, &groups); err != nil {
		return nil, fmt.Errorf("failed to get all groups: %w", err)
	}

	return groups, nil
}

func (c *Client) GetAllEmployees() ([]models.EmployeeListItem, error) {
	url := fmt.Sprintf("%s/employees/all", c.baseURL)

	var employees []models.EmployeeListItem
	if err := c.makeRequest(url, &employees); err != nil {
		return nil, fmt.Errorf("failed to get all employees: %w", err)
	}

	return employees, nil
}

func (c *Client) GetAllFaculties() ([]models.Faculty, error) {
	url := fmt.Sprintf("%s/faculties", c.baseURL)

	var faculties []models.Faculty
	if err := c.makeRequest(url, &faculties); err != nil {
		return nil, fmt.Errorf("failed to get all faculties: %w", err)
	}

	return faculties, nil
}

func (c *Client) GetAllDepartments() ([]models.Department, error) {
	url := fmt.Sprintf("%s/departments", c.baseURL)

	var departments []models.Department
	if err := c.makeRequest(url, &departments); err != nil {
		return nil, fmt.Errorf("failed to get all departments: %w", err)
	}

	return departments, nil
}

func (c *Client) GetAllSpecialities() ([]models.Speciality, error) {
	url := fmt.Sprintf("%s/specialities", c.baseURL)

	var specialities []models.Speciality
	if err := c.makeRequest(url, &specialities); err != nil {
		return nil, fmt.Errorf("failed to get all specialities: %w", err)
	}

	return specialities, nil
}

func (c *Client) GetEmployeeAnnouncements(urlID string) ([]models.Announcement, error) {
	url := fmt.Sprintf("%s/announcements/employees?url-id=%s", c.baseURL, urlID)

	var announcements []models.Announcement
	if err := c.makeRequest(url, &announcements); err != nil {
		return nil, fmt.Errorf("failed to get employee announcements: %w", err)
	}

	return announcements, nil
}

func (c *Client) GetDepartmentAnnouncements(departmentID int) ([]models.Announcement, error) {
	url := fmt.Sprintf("%s/announcements/departments?id=%d", c.baseURL, departmentID)

	var announcements []models.Announcement
	if err := c.makeRequest(url, &announcements); err != nil {
		return nil, fmt.Errorf("failed to get department announcements: %w", err)
	}

	return announcements, nil
}

func (c *Client) GetAllAuditories() ([]models.Auditory, error) {
	url := fmt.Sprintf("%s/auditories", c.baseURL)

	var auditories []models.Auditory
	if err := c.makeRequest(url, &auditories); err != nil {
		return nil, fmt.Errorf("failed to get all auditories: %w", err)
	}

	return auditories, nil
}

func (c *Client) GetGroupLastUpdateDate(groupNumber string) (*models.LastUpdateDate, error) {
	url := fmt.Sprintf("%s/last-update-date/student-group?groupNumber=%s", c.baseURL, groupNumber)

	var updateDate models.LastUpdateDate
	if err := c.makeRequest(url, &updateDate); err != nil {
		return nil, fmt.Errorf("failed to get group last update date: %w", err)
	}

	return &updateDate, nil
}

func (c *Client) GetGroupLastUpdateDateByID(groupID int) (*models.LastUpdateDate, error) {
	url := fmt.Sprintf("%s/last-update-date/student-group?id=%d", c.baseURL, groupID)

	var updateDate models.LastUpdateDate
	if err := c.makeRequest(url, &updateDate); err != nil {
		return nil, fmt.Errorf("failed to get group last update date by ID: %w", err)
	}

	return &updateDate, nil
}

func (c *Client) GetEmployeeLastUpdateDate(urlID string) (*models.LastUpdateDate, error) {
	url := fmt.Sprintf("%s/last-update-date/employee?url-id=%s", c.baseURL, urlID)

	var updateDate models.LastUpdateDate
	if err := c.makeRequest(url, &updateDate); err != nil {
		return nil, fmt.Errorf("failed to get employee last update date: %w", err)
	}

	return &updateDate, nil
}

func (c *Client) GetEmployeeLastUpdateDateByID(employeeID int) (*models.LastUpdateDate, error) {
	url := fmt.Sprintf("%s/last-update-date/employee?id=%d", c.baseURL, employeeID)

	var updateDate models.LastUpdateDate
	if err := c.makeRequest(url, &updateDate); err != nil {
		return nil, fmt.Errorf("failed to get employee last update date by ID: %w", err)
	}

	return &updateDate, nil
}

func (c *Client) GetCurrentWeek() (int, error) {
	url := fmt.Sprintf("%s/schedule/current-week", c.baseURL)

	var week int
	if err := c.makeRequest(url, &week); err != nil {
		return 0, fmt.Errorf("failed to get current week: %w", err)
	}

	return week, nil
}

func (c *Client) makeRequest(url string, target interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
