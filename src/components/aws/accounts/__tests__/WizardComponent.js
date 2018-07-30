import React from 'react';
import WizardComponent, {
  StepOne,
  StepTwo,
  StepThree
} from '../WizardComponent';
import Dialog, {
  DialogContent
} from 'material-ui/Dialog';
import Stepper, {
  Step,
  StepButton
} from 'material-ui/Stepper';
import Spinner from 'react-spinkit';
import Misc from '../../../misc';
import { shallow } from 'enzyme';
import Input from "react-validation/build/input";
import Button from "react-validation/build/button";
import Form from "react-validation/build/form";

const Picture = Misc.Picture;

const external = {
  external: "external",
  accountId: "accountId"
};

const accountEmpty = {
  status: true,
  value: null
};

const account = {
  ...accountEmpty,
  value: {
    id: 42,
    roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE",
    pretty: "pretty"
  }
};

const accountWaiting = {
  status: false
};

const accountError = {
  status: true,
  error: Error("Error")
};

const billEmpty = {
  status: true,
  value: null
};

const bill = {
  ...billEmpty,
  value: {
    id: 42,
    roleArn: "arn:aws:iam::000000000000:role/TEST_ROLE",
    pretty: "pretty"
  }
};

const billWaiting = {
  status: false
};

const billError = {
  status: true,
  error: Error("Error")
};

describe('<WizardComponent />', () => {

  const props = {
    submitAccount: jest.fn(),
    clearAccount: jest.fn(),
    submitBucket: jest.fn(),
    clearBucket: jest.fn()
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <WizardComponent /> component', () => {
    const wrapper = shallow(<WizardComponent {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <Dialog /> component', () => {
    const wrapper = shallow(<WizardComponent {...props}/>);
    const children = wrapper.find(Dialog);
    expect(children.length).toBe(1);
  });

  it('renders a <DialogContent /> component', () => {
    const wrapper = shallow(<WizardComponent {...props}/>);
    const children = wrapper.find(DialogContent);
    expect(children.length).toBe(1);
  });

  it('can open and close dialog', () => {
    const wrapper = shallow(<WizardComponent {...props}/>);
    expect(wrapper.state('open')).toBe(false);
    expect(props.clearAccount).not.toHaveBeenCalled();
    expect(props.clearBucket).not.toHaveBeenCalled();
    wrapper.instance().openDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(true);
    expect(props.clearAccount).toHaveBeenCalledTimes(1);
    expect(props.clearBucket).toHaveBeenCalledTimes(1);
    wrapper.instance().closeDialog({ preventDefault(){} });
    expect(wrapper.state('open')).toBe(false);
    expect(props.clearAccount).toHaveBeenCalledTimes(2);
    expect(props.clearBucket).toHaveBeenCalledTimes(2);
    wrapper.instance().closeDialog();
    expect(wrapper.state('open')).toBe(false);
    expect(props.clearAccount).toHaveBeenCalledTimes(3);
    expect(props.clearBucket).toHaveBeenCalledTimes(3);
  });

  it('renders a <Stepper /> component', () => {
    const wrapper = shallow(<WizardComponent {...props}/>);
    const children = wrapper.find(Stepper);
    expect(children.length).toBe(1);
  });

  it('renders three <Step /> components', () => {
    const wrapper = shallow(<WizardComponent {...props}/>);
    const children = wrapper.find(Step);
    expect(children.length).toBe(3);
  });

  it('renders three <StepButton /> components', () => {
    const wrapper = shallow(<WizardComponent {...props}/>);
    const children = wrapper.find(StepButton);
    expect(children.length).toBe(3);
  });

  it('can go to next and previous step', () => {
    const wrapper = shallow(<WizardComponent {...props}/>);
    expect(wrapper.state('activeStep')).toBe(0);
    wrapper.instance().nextStep();
    expect(wrapper.state('activeStep')).toBe(1);
    wrapper.instance().previousStep();
    expect(wrapper.state('activeStep')).toBe(0);
  });

});

describe('<StepOne />', () => {

  const props = {
    external,
    next: jest.fn(),
    close: jest.fn()
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <StepOne /> component', () => {
    const wrapper = shallow(<StepOne {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <div /> component for tutorial', () => {
    const wrapper = shallow(<StepOne {...props}/>);
    const children = wrapper.find("div.tutorial");
    expect(children.length).toBe(1);
  });

  it('renders a <Picture /> component in <div /> tutorial', () => {
    const wrapper = shallow(<StepOne {...props}/>);
    const picture = wrapper.find(Picture);
    expect(picture.length).toBe(1);
  });

  it('renders a <Form /> component', () => {
    const wrapper = shallow(<StepOne {...props}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('renders 1 <Button /> component in <Form />', () => {
    const wrapper = shallow(<StepOne {...props}/>);
    const form = wrapper.find(Form);
    const button = form.find(Button);
    expect(button.length).toBe(1);
  });

  it('renders 1 <button /> component in <Form />', () => {
    const wrapper = shallow(<StepOne {...props}/>);
    const form = wrapper.find(Form);
    const button = form.find("button");
    expect(button.length).toBe(1);
  });

  it('can submit', () => {
    const wrapper = shallow(<StepOne {...props}/>);
    expect(props.next).not.toHaveBeenCalled();
    wrapper.instance().submit({ preventDefault() {} });
    expect(props.next).toHaveBeenCalled();
  });

});

describe('<StepTwo />', () => {

  const props = {
    account: accountEmpty,
    external,
    next: jest.fn(),
    back: jest.fn(),
    submit: jest.fn(),
    close: jest.fn()
  };

  const propsWaiting = {
    ...props,
    account: accountWaiting
  };

  const propsDone = {
    ...props,
    account
  };

  const propsError = {
    ...props,
    account: accountError
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <StepTwo /> component', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <div /> component for tutorial', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    const children = wrapper.find("div.tutorial");
    expect(children.length).toBe(1);
  });

  it('renders a <Picture /> component in <div /> tutorial', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    const picture = wrapper.find(Picture);
    expect(picture.length).toBe(1);
  });

  it('renders a <Form /> component', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('renders 2 <Input /> components in <Form />', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    const form = wrapper.find(Form);
    const inputs = form.find(Input);
    expect(inputs.length).toBe(2);
  });

  it('renders 1 <Button /> component in <Form />', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    const form = wrapper.find(Form);
    const button = form.find(Button);
    expect(button.length).toBe(1);
  });

  it('renders 2 <button /> components in <Form />', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    const form = wrapper.find(Form);
    const button = form.find("div.btn");
    expect(button.length).toBe(2);
  });

  it('renders a <Spinner /> component if waiting for response', () => {
    let wrapper = shallow(<StepTwo {...props}/>);
    let spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(0);
    wrapper = shallow(<StepTwo {...propsWaiting}/>);
    spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(1);
  });

  it('renders an alert if there is an error', () => {
    const wrapper = shallow(<StepTwo {...propsError}/>);
    const error = wrapper.find("div.alert");
    expect(error.length).toBe(1);
  });

  it('can submit', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    const instance = wrapper.instance();
    instance.form = {
      getValues: () => ({
        roleArn: "roleArn",
        pretty: "pretty"
      })
    };
    expect(props.submit).not.toHaveBeenCalled();
    wrapper.instance().submit({ preventDefault() {} });
    expect(props.submit).toHaveBeenCalled();
  });

  it('can go to next page when data is available', () => {
    const wrapper = shallow(<StepTwo {...props}/>);
    expect(props.next).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(propsWaiting);
    expect(props.next).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(propsDone);
    expect(props.next).toHaveBeenCalled();
  });

});

describe('<StepThree />', () => {

  const props = {
    account,
    external,
    submit: jest.fn(),
    close: jest.fn()
  };

  const propsWithoutBill = {
    ...props,
    bill: null
  };

  const propsWithoutAccount = {
    ...props,
    account: null,
    bill
  };

  const propsWaiting = {
    ...props,
    bill: billWaiting
  };

  const propsDone = {
    ...props,
    bill
  };

  const propsError = {
    ...props,
    bill: billError
  };

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it('renders a <StepThree /> component', () => {
    const wrapper = shallow(<StepThree {...propsWithoutAccount}/>);
    expect(wrapper.length).toBe(1);
  });

  it('renders a <div /> component for tutorial', () => {
    const wrapper = shallow(<StepThree {...propsWithoutAccount}/>);
    const children = wrapper.find("div.tutorial");
    expect(children.length).toBe(1);
  });

  it('renders a <Form /> component', () => {
    const wrapper = shallow(<StepThree {...propsWithoutAccount}/>);
    const form = wrapper.find(Form);
    expect(form.length).toBe(1);
  });

  it('renders 2 <Input /> component in <Form />', () => {
    const wrapper = shallow(<StepThree {...propsWithoutAccount}/>);
    const form = wrapper.find(Form);
    const inputs = form.find(Input);
    expect(inputs.length).toBe(2);
  });

  it('renders 1 <Button /> component in <Form />', () => {
    let wrapper = shallow(<StepThree {...propsWithoutAccount}/>);
    let form = wrapper.find(Form);
    let button = form.find(Button);
    expect(button.length).toBe(1);
    expect(button.prop("disabled")).toBe(true);
    wrapper = shallow(<StepThree {...propsWaiting}/>);
    form = wrapper.find(Form);
    button = form.find(Button);
    expect(button.length).toBe(1);
    expect(button.prop("disabled")).toBe(false);
  });

  it('renders 1 <button /> component in <Form />', () => {
    const wrapper = shallow(<StepThree {...propsWithoutAccount}/>);
    const form = wrapper.find(Form);
    const button = form.find("div.btn");
    expect(button.length).toBe(1);
  });

  it('renders a <Spinner /> component if waiting for response', () => {
    let wrapper = shallow(<StepThree {...propsWithoutAccount}/>);
    let spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(0);
    wrapper = shallow(<StepThree {...propsWaiting}/>);
    spinner = wrapper.find(Spinner);
    expect(spinner.length).toBe(2);
  });

  it('renders an alert if there is an error', () => {
    const wrapper = shallow(<StepThree {...propsError}/>);
    const error = wrapper.find("div.alert");
    expect(error.length).toBe(1);
  });

  it('can submit', () => {
    const wrapper = shallow(<StepThree {...propsWithoutBill}/>);
    const instance = wrapper.instance();
    instance.form = {
      getValues: () => ({
        bucket: "s3://account/path/to/bills"
      })
    };
    expect(props.submit).not.toHaveBeenCalled();
    wrapper.instance().submit({ preventDefault() {} });
    expect(props.submit).toHaveBeenCalled();
  });

  it('can go to next page when data is available', () => {
    const wrapper = shallow(<StepThree {...propsWithoutBill}/>);
    expect(props.close).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(propsWaiting);
    expect(props.close).not.toHaveBeenCalled();
    wrapper.instance().componentWillReceiveProps(propsDone);
    expect(props.close).toHaveBeenCalled();
  });

});
