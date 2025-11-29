import type {
  ScheduleResponse,
  StudentGroupListItem,
  EmployeeListItem,
  ApiError,
  RefreshResponse,
} from '@/types/api';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({
        error: `HTTP ${response.status}: ${response.statusText}`,
      }));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // Schedule endpoints
  async getGroupSchedule(
    groupNumber: string,
    useCache: boolean = true
  ): Promise<ScheduleResponse> {
    return this.request<ScheduleResponse>(
      `/schedule/group/${groupNumber}?useCache=${useCache}`
    );
  }

  async getEmployeeSchedule(
    urlId: string,
    useCache: boolean = true
  ): Promise<ScheduleResponse> {
    return this.request<ScheduleResponse>(
      `/schedule/employee/${urlId}?useCache=${useCache}`
    );
  }

  async refreshGroupSchedule(groupNumber: string): Promise<RefreshResponse> {
    return this.request<RefreshResponse>(`/schedule/group/${groupNumber}/refresh`, {
      method: 'POST',
    });
  }

  async refreshEmployeeSchedule(urlId: string): Promise<RefreshResponse> {
    return this.request<RefreshResponse>(`/schedule/employee/${urlId}/refresh`, {
      method: 'POST',
    });
  }

  // Groups endpoints
  async getAllGroups(useCache: boolean = true): Promise<StudentGroupListItem[]> {
    return this.request<StudentGroupListItem[]>(
      `/groups?useCache=${useCache}`
    );
  }

  async getGroupByNumber(groupNumber: string): Promise<StudentGroupListItem> {
    return this.request<StudentGroupListItem>(`/groups/${groupNumber}`);
  }

  async refreshGroups(): Promise<RefreshResponse> {
    return this.request<RefreshResponse>('/groups/refresh', {
      method: 'POST',
    });
  }

  // Employees endpoints
  async getAllEmployees(useCache: boolean = true): Promise<EmployeeListItem[]> {
    return this.request<EmployeeListItem[]>(
      `/employees?useCache=${useCache}`
    );
  }

  async getEmployeeByUrlId(urlId: string): Promise<EmployeeListItem> {
    return this.request<EmployeeListItem>(`/employees/${urlId}`);
  }

  async refreshEmployees(): Promise<RefreshResponse> {
    return this.request<RefreshResponse>('/employees/refresh', {
      method: 'POST',
    });
  }
}

export const apiClient = new ApiClient();

