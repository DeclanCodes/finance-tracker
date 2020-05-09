import React from "react";
import api from "./api";
import CreateCategoryForm from "./CreateCategoryForm";
import CategoryRow from "./CategoryRow";

class CategoryPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      categories: []
    };
    this.handleCreate = this.handleCreate.bind(this);
    this.handleUpdate = this.handleUpdate.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
  }

  isAccountCategory() {
    return this.props.categoryType === "Account"
  }

  getCategories(isAccountCategory) {
    return isAccountCategory
        ? api.getAccountCategories()
        : api.getExpenseCategories()
  }

  handleCreate(values) {
    const isAccountCategory = this.isAccountCategory()
    const p = isAccountCategory
      ? api.createAccountCategory(values)
      : api.createExpenseCategory(values)

    p.then(() => this.getCategories(isAccountCategory))
      .then(response => this.setState({ categories: response.data }))
  }

  handleDelete(uuid) {
    const isAccountCategory = this.isAccountCategory()
    const p = isAccountCategory
      ? api.deleteAccountCategory(uuid)
      : api.deleteExpenseCategory(uuid)

    p.then(() => this.getCategories(isAccountCategory))
      .then(response => this.setState({ categories: response.data }))
  }

  handleUpdate(values) {
    const isAccountCategory = this.isAccountCategory()
    const p = isAccountCategory
      ? api.updateAccountCategory(values)
      : api.updateExpenseCategory(values)

    p.then(() => this.getCategories(isAccountCategory))
      .then(response => this.setState({ categories: response.data }))
  }

  componentDidMount() {
    this.getCategories(this.isAccountCategory())
      .then(response => response.data)
      .then(data => this.setState({ categories: data }))
  }

  render() {
    return (
      <div>
        <h1>{this.props.categoryType} Categories</h1>
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {this.state.categories.length > 0 ? (
              this.state.categories.map(category => (
                (
                  <CategoryRow
                    key={category.uuid}
                    categoryType={this.props.categoryType}
                    category={category}
                    handleUpdate={this.handleUpdate}
                    handleDelete={this.handleDelete}
                  />
                )
              ))
            ) : (
              <tr>
                <td colSpan={3}>No {this.props.categoryType} Categories</td>
              </tr>
            )}
          </tbody>
        </table>
        <CreateCategoryForm
          categoryType={this.props.categoryType}
          doSubmit={this.handleCreate}
        />
      </div>
    );
  }
}

export default CategoryPage;