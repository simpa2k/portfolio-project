import React from 'react';
import './App.css';
import ValidatedForm from './ValidatedForm';
import { Card, Container, Row, Col, Form, InputGroup, Button } from 'react-bootstrap';

function getMessages(setter) {
  fetch('/api/message')
    .then((res) => res.json())
    .then((data) => {
      console.log('Received messages from server:', data);
      const mapped = data.map((d) => Object.entries(d).map(([key, value]) => <p>{key}: {value}</p>))

      console.log(mapped)

      setter(data)
    });
}

function App() {
  const [data, setData] = React.useState([]);
  const [message, setMessage] = React.useState({})
  const [customFieldName, setCustomFieldName] = React.useState("");
  const [customFieldType, setCustomFieldType] = React.useState(null);
  const [customFieldRequired, setCustomFieldRequired] = React.useState(false);

  let initialFormGroups = [
    <Form.Group>
      <Form.Label>Message</Form.Label>
      <InputGroup>
        <Form.Control as="textarea" rows="4" placeholder="Write your message here" onChange={(e) => setMessage({...message, 'message': e.target.value})} required/>
        <Form.Control.Feedback type="invalid">
          Please write a message.
        </Form.Control.Feedback>
      </InputGroup>
    </Form.Group>
  ]

  const [formGroups, setFormGroups] = React.useState(initialFormGroups)

  React.useEffect(() => getMessages(setData), []);

  const postMessages = () => {
    const postBody = JSON.stringify(message);
    console.log(postBody);

    fetch('/api/message', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: postBody
    }).then(() => getMessages(setData)); // TODO: Check response
  };

  const deleteFormField = (i, key) => {
    let copiedFormGroups = [...formGroups];

    // Remove form group
    copiedFormGroups.splice(i, 1);
    setFormGroups(copiedFormGroups);

    // Remove key from message
    setMessage(message => {
      delete message[key];
      return message;
    });
  };
  
  const addNewFormField = () => {
    let attributes = {};
    if (customFieldRequired) {
      attributes['required'] = 'required'
    }

    if (customFieldType === 'textarea') {
      attributes['as'] = customFieldType
    } else {
      attributes['type'] = customFieldType
    }

    setFormGroups([...formGroups,
      <Form.Group>
        <Form.Label>{customFieldName}</Form.Label>
        <InputGroup>
          <Form.Control {...attributes} onChange={(e => setMessage(message => { 
            return {...message, [customFieldName]: e.target.value}
          }))}/>
          <Button variant="danger" onClick={() => deleteFormField(formGroups.length + 1, customFieldName)}>-</Button>
        </InputGroup>
      </Form.Group>
    ])

    setCustomFieldName("");
  };

  return (
    <Container>
        {!data ? <p>Loading</p> : data.map((d, i) =>
          <Row key={i}>
            <Col xs={{span: 4, offset: 4}}>
              <Card>
                <Card.Body>
                  <Card.Title>Message</Card.Title>
                  <Card.Text>{Object.entries(d).map(([key, value]) => <p>{key}: {value}</p>)}</Card.Text>
                </Card.Body>
              </Card>
            </Col>
          </Row>
        )}
      <Row>
        <Col xs={{span: 6, offset: 3}}>
          <ValidatedForm onSubmit={postMessages}>
            {formGroups}
          </ValidatedForm>
        </Col>
      </Row>
      <Row>
        <Col xs={{span: 6, offset: 3}}>
          <ValidatedForm onSubmit={addNewFormField}>
            <Form.Group>
              <Form.Label>Name</Form.Label>
              <InputGroup>
                <Form.Control type="input" placeholder="Name" value={customFieldName} onChange={e => setCustomFieldName(e.target.value)} required/>
                <Form.Control.Feedback type="invalid">You need to provide a name for the new field</Form.Control.Feedback>
              </InputGroup>
            </Form.Group>
            <Form.Group>
              <Form.Label>Type</Form.Label>
              <InputGroup>
                <Form.Select aria-label="Select custom field type" onChange={e => setCustomFieldType(e.target.value)}>
                  <option value="input">Input</option>
                  <option value="textarea">Textarea</option>
                </Form.Select>
                <Form.Control.Feedback type="invalid">You need to provide a type for the new field</Form.Control.Feedback>
              </InputGroup>
            </Form.Group>
            <Form.Group>
              <Form.Check label="Required" onChange={e => setCustomFieldRequired(!customFieldRequired)}/>
            </Form.Group>
          </ValidatedForm>
        </Col>
      </Row>
    </Container>
  );
}

export default App;
