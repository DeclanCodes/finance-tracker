import React from 'react';
import { Button } from '../Button';
import { EntityForm } from '../forms/EntityForm';
import { ModifyRowPanel } from './ModifyRowPanel';
import { helpers } from '../../../common/helpers';

class EntityRow extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isEditing: false
    };
    this.stopEditing = this.stopEditing.bind(this);
  }

  setEditing(val) {
    this.setState({ isEditing: val });
  }

  stopEditing() {
    this.setEditing(false);
  }

  render() {
    const e = this.props.entity;

    return (
      <tr>
        {e.hasOwnProperty('name') && <td>{e.name}</td>}
        {e.hasOwnProperty('account') && <td>{e.account.name}</td>}
        {e.hasOwnProperty('category') && <td>{e.category.name}</td>}
        {e.hasOwnProperty('description') && <td>{e.description}</td>}
        {e.hasOwnProperty('date') && <td>{helpers.displayDate(e.date)}</td>}
        {e.hasOwnProperty('amount') && <td>{`$${e.amount}`}</td>}
        <td>
          {this.state.isEditing ? (
            <div>
              <EntityForm
                entityName={this.props.entityName}
                entity={e}
                getInitialValues={this.props.getInitialValues}
                isCreateMode={false}
                options={this.props.options}
                doExtraModifications={this.props.doExtraModifications}
                doSubmit={this.props.handleUpdate}
                doFinalState={this.stopEditing}
              />
              <Button
                name='Cancel'
                handleFunc={this.stopEditing}
              />
            </div>
          ) : (
            <ModifyRowPanel
              handleEdit={() => this.setEditing(true)}
              handleDelete={() => this.props.handleDelete(e.uuid)}
            />
          )}
        </td>
      </tr>
    );
  }
}

export default EntityRow;
