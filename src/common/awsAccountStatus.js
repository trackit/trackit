import React from "react";

const errors = {
  "error": 0,
  "not_started": 1,
  "in_progress": 2,
  "ok": 3
};

const getBadge = (status) => {
  switch (status.value) {
    case "error":
      return (<i className="fa account-badge fa-times-circle"/>);
    case "not_started":
    case "in_progress":
      return (<i className="fa account-badge fa-clock-o"/>);
    case "ok":
      return (<i className="fa account-badge fa-check-circle"/>);
    default:
      return null;
  }
};

const getInformationBanner = (status) => {
  switch (status.value) {
    case "error":
      return (<div className="alert alert-danger account-badge-information-banner">Import failed, please check your bills locations.</div>);
    case "not_started":
    case "in_progress":
      return (<div className="alert alert-warning account-badge-information-banner">Import may take up to 24 hours, please wait.</div>);
    case "ok":
    default:
      return null;
  }
};

const getAWSAccountStatus = (account) => {
  const status = {
    value: "ok",
    details: []
  };
  if (account.status) {
    status.value = account.status.value;
    status.details.push("Account : " + account.status.detail);
  }
  let hasBillRepositoryError = false;
  getBillRepositoriesStatuses(account.billRepositories).forEach((billRepository) => {
    if (billRepository) {
      if (errors[status.value] > errors[billRepository.value])
        status.value = billRepository.value;
      if (billRepository.value === "error" && !hasBillRepositoryError) {
        status.details.push("Bill locations : See in dedicated section for more details");
        hasBillRepositoryError = true;
      }
    }
  });
  return status;
};

const getBillRepositoriesStatuses = (billRepositories) => {
  if (billRepositories)
    return billRepositories.map((billRepository) => (billRepository.status));
  return [];
};

export default {
  getAWSAccountStatus,
  getBadge,
  getInformationBanner
};
