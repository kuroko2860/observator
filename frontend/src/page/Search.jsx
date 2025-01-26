import { FormProvider, useForm } from "react-hook-form";
import CustomContainer from "../component/CustomContainer";
import { ServiceNameInput } from "../component/Input";
import { SubmitButtons } from "../component/Common";
import { Grid2 } from "@mui/material";
import { useNavigate } from "react-router-dom";

const Search = () => {
  const methods = useForm({});
  const navigate = useNavigate();
  const onSubmit = (data) => {
    navigate(`/service-detail/${data.service_name}`);
  };
  return (
    <CustomContainer>
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(onSubmit)}>
          <Grid2 container spacing={2}>
            <ServiceNameInput />
            <SubmitButtons />
          </Grid2>
        </form>
      </FormProvider>
    </CustomContainer>
  );
};

export default Search;
