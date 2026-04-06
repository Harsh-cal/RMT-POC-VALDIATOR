import axios from "axios";

const BASE_URL = "http://localhost:8080/api/dev/v1";

export const getReleasesHistory = async ({ limit = 50, offset = 0 } = {}) => {
  const response = await axios.get(`${BASE_URL}/releases/history`, {
    params: { limit, offset },
  });
  return response.data;
};

export const getTrendsData = async ({ days = 90 } = {}) => {
  const response = await axios.get(`${BASE_URL}/releases/trends`, {
    params: { days },
  });
  return response.data;
};

export const getRecurringIssues = async ({ days = 90 } = {}) => {
  const response = await axios.get(`${BASE_URL}/issues/recurring`, {
    params: { days },
  });
  return response.data;
};

export const compareReleases = async (id1, id2) => {
  const response = await axios.get(`${BASE_URL}/releases/${id1}/compare/${id2}`);
  return response.data;
};
