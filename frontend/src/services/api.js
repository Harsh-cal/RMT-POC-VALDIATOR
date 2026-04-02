import axios from "axios";

const BASE_URL = "http://localhost:8080/api/dev/v1";

export const validateRelease = async (payload) => {
  const response = await axios.post(`${BASE_URL}/validate`, payload);
  console.log(response);
  
  return response.data;
};

export const askValidationBot = async ({ question, result, history = [] }) => {
  const response = await axios.post(`${BASE_URL}/validate/chat`, {
    question,
    result,
    history,
  });
  return response.data;
};

export const exportValidationReport = async ({ releaseId, releaseName, format }) => {
  console.log("Exporting report for", { releaseId, releaseName, format });
  const response = await axios.post(
    `${BASE_URL}/validate/export`,
    {
      release_id: releaseId || "",
      release_name: releaseName || "",
      format,
    },
    {
      responseType: "blob",
    }
  );

  const contentDisposition = response.headers["content-disposition"] || "";
  const fileNameMatch = contentDisposition.match(/filename=([^;]+)/i);
  const rawFileName = fileNameMatch?.[1]?.trim() || `validation-report.${format}`;
  const fileName = rawFileName.replace(/^"|"$/g, "");

  const blob = new Blob([response.data], {
    type: response.headers["content-type"] || "application/octet-stream",
  });

  const url = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = fileName;
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);

  console.log(link)
};