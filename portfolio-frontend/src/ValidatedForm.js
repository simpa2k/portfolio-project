import React from 'react';
import { Form, Button } from 'react-bootstrap';

function ValidatedForm(props) {
  const [validated, setValidated] = React.useState(false)

  const handleSubmit = (event) => {
    const form = event.currentTarget;

    event.preventDefault();
    event.stopPropagation();

    if (form.checkValidity()) {
      setValidated(true);
      props.onSubmit();
    }
  };

  return (
    <Form noValidate validated={validated} onSubmit={handleSubmit}>
      {props.children}
      <Button type="submit">Send</Button>
    </Form>
  )
}

export default ValidatedForm
